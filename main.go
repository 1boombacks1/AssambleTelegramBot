package main

import (
	"cryptoWallet/models"
	"encoding/json"
	"flag"
	"fmt"
	"io/fs"
	"io/ioutil"
	"log"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"cryptoWallet/keyboards"
	"cryptoWallet/messages"
)

func main() {
	bot, err := tgbotapi.NewBotAPI(mustToken())
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = false

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")

		if update.Message.IsCommand() {
			switch update.Message.Text {
			case "/start":
				chatID := update.Message.Chat.ID
				username := update.Message.From.UserName

				if userIsExist(chatID) {
					errMsg := tgbotapi.NewMessage(chatID, messages.UserExistMessage)
					bot.Send(errMsg)
					continue
				}

				if err := addUser(chatID, username); err != nil {
					errMsg := tgbotapi.NewMessage(chatID, messages.AddToDbError)
					bot.Send(errMsg)
					continue
				}

				log.Printf("Registred user: %s", username)

				msg.Text = messages.StartMessage
				msg.ReplyMarkup = keyboards.AssembleKeyboard

			case "/help":
				msg.Text = messages.HelpMessage

			case "/delete":
				chatID := update.Message.Chat.ID

				if !userIsExist(chatID) {
					msg.Text = messages.UserNotExistMessage
					bot.Send(msg)
					continue
				}

				if err := deleteUser(chatID); err != nil {
					log.Printf("Ошибка удаление с бд: %s", err.Error())
					msg.Text = messages.RemoveToDbError
					bot.Send(msg)
					continue
				}

				log.Printf("Deleted user: %s", update.Message.Chat.UserName)

				msg.Text = messages.RemoveToDbMessage
				msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
			}
		}

		if update.Message.Text == keyboards.AssembleText {
			users, err := getUsers()
			if err != nil {
				log.Printf("Ошибка получения пользователей: %s", err.Error())
			}

			if dur, ok := checkForSending(update.Message.Chat.ID); !ok {
				log.Printf("%s trying gather to party", update.Message.Chat.UserName)
				text := fmt.Sprintf("Вы сможете подтянуть рыцарей через: %s ⌛", dur)
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, text)
				bot.Send(msg)
				continue
			}

			convenerName := update.Message.From.UserName
			text := fmt.Sprintf("@%s созывает на стык!🤘", convenerName)

			for _, user := range users {
				msgForUser := tgbotapi.NewMessage(user.ChatID, text)
				if _, err := bot.Send(msgForUser); err != nil {
					errMsg := tgbotapi.NewMessage(update.Message.Chat.ID, messages.AssembleError)
					bot.Send(errMsg)
					break
				}
			}

			log.Printf("%s called everyone together", update.Message.Chat.UserName)
			updateAttempts(update.Message.Chat.ID)
			continue

		} else {
			if msg.Text == "" {
				log.Printf("%s writed: %s", update.Message.Chat.UserName, update.Message.Text)
				msg.Text = messages.NotCorrectCommandMessage
			}
			if _, err := bot.Send(msg); err != nil {
				log.Printf("Ошибка отправки: %s", err.Error())
				continue
			}
		}
	}
}

func checkForSending(chatID int64) (string, bool) {
	users, _ := getUsers()
	for _, user := range users {
		if chatID == user.ChatID {
			duration := time.Since(user.Timer).Round(time.Second)
			if user.Attempts > 0 {
				return "", true
			} else if user.Attempts == 0 &&
				duration.Seconds() < (time.Minute*5).Seconds() {
				return time.Until(user.Timer.Add(5 * time.Minute)).Round(time.Second).String(), false
			}
			break
		}
	}
	return "", true
}

func updateAttempts(chatID int64) {
	users, _ := getUsers()
	for i, user := range users {
		if user.ChatID == chatID {
			if user.Attempts == 2 {
				users[i].Attempts--
				users[i].Timer = time.Now()
			} else if user.Attempts == 1 {
				users[i].Attempts--
			} else {
				duration := time.Since(user.Timer)
				if duration.Seconds() > (time.Minute * 5).Seconds() {
					users[i].Attempts = 2
					users[i].Timer = time.Time{}
				}
			}
			break
		}
	}
	writeUsers(users)
}

func userIsExist(chatID int64) bool {
	users, err := getUsers()
	if err != nil {
		return false
	}
	for _, user := range users {
		if user.ChatID == chatID {
			return true
		}
	}
	return false
}

func deleteUser(chatID int64) error {
	users, err := getUsers()
	if err != nil {
		return err
	}

	updateUsers := users

	for i, user := range users {
		if user.ChatID == chatID {
			updateUsers = removeIndex(users, i)
			break
		}
	}

	if err := writeUsers(updateUsers); err != nil {
		return err
	}

	return nil
}

func addUser(chatID int64, username string) error {
	users, err := getUsers()
	if err != nil {
		return err
	}

	user := models.User{
		ChatID:   chatID,
		Username: username,
		Attempts: 2,
		Timer:    time.Time{},
	}
	users = append(users, user)

	if err := writeUsers(users); err != nil {
		return err
	}
	return nil
}

func writeUsers(users []models.User) error {
	data, err := json.Marshal(users)
	if err != nil {
		log.Printf("Не удалось маршануть данные: %s", err.Error())
		return err
	}

	if err = ioutil.WriteFile("data/users.json", data, fs.ModePerm); err != nil {
		log.Printf("Не удалось записать данные в файл: %s", err.Error())
		return err
	}
	return nil
}

func removeIndex(users []models.User, i int) []models.User {
	newArr := make([]models.User, 0)
	newArr = append(newArr, users[:i]...)
	return append(newArr, users[i+1:]...)
}

func getUsers() ([]models.User, error) {
	file, err := ioutil.ReadFile("data/users.json")
	if err != nil {
		log.Printf("Не удалось прочитать файл: users.json <%s>", err.Error())
		return nil, err
	}

	var users []models.User

	if err = json.Unmarshal(file, &users); err != nil {
		log.Printf("Не удалось Анмаршануть данные: %s", err.Error())
		return nil, err
	}

	return users, nil
}

func mustToken() string {
	token := flag.String(
		"bot-token",
		"",
		"telegram bot token for the application to work",
	)

	flag.Parse()

	if *token == "" {
		log.Fatal("No token entered")
	}
	return *token
}
