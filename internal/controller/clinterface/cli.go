//nolint:gocognit
package clinterface

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/AntonNikol/anti-bruteforce/internal/controller/httpapi/handlers"
	"github.com/AntonNikol/anti-bruteforce/internal/domain/entity"
	"github.com/AntonNikol/anti-bruteforce/internal/domain/service"
	"github.com/c-bata/go-prompt"
)

var suggestions = []prompt.Suggest{
	{Text: "blacklist add [ip_address] [mask]", Description: "Add ip net to blacklist"},
	{Text: "blacklist remove [ip_address] [mask]", Description: "Remove ip net to blacklist"},
	{Text: "blacklist get", Description: "Get ip list from blacklist"},
	{Text: "whitelist add [ip_address] [mask]", Description: "Add ip net to whitelist"},
	{Text: "whitelist remove [ip_address] [mask]", Description: "Remove ip net to whitelist"},
	{Text: "whitelist get", Description: "Get ip list from whitelist"},
	{Text: "bucket remove [login] [ip_address]", Description: "Remove login and ip address from bucket"},
	{Text: "help", Description: "Display list of commands"},
	{Text: "exit", Description: "Exit anti bruteforce app"},
}

type CommandLineInterface struct {
	serviceAuth *service.Authorization
	serviceWL   *service.WhiteList
	serviceBL   *service.BlackList
}

// New создает экземпляр командного интерфейса CommandLineInterface и инициализирует все необходимые сервисы.
func New(
	serviceAuth *service.Authorization,
	serviceWL *service.WhiteList,
	serviceBL *service.BlackList,
) *CommandLineInterface {
	return &CommandLineInterface{serviceAuth: serviceAuth, serviceWL: serviceWL, serviceBL: serviceBL}
}

// Run запускает командный интерфейс и обрабатывает ввод пользователя.
func (c *CommandLineInterface) Run(ch chan os.Signal) {
	executor := prompt.Executor(func(s string) {
		s = strings.TrimSpace(s)
		setCommand := strings.Split(s, " ")
		switch setCommand[0] {
		case "blacklist":
			switch setCommand[1] {
			case "add":
				if len(setCommand) != 4 {
					break
				}
				c.addIPToBl(entity.IPNetwork{
					IP:   setCommand[2],
					Mask: setCommand[3],
				})
			case "remove":
				if len(setCommand) != 4 {
					break
				}
				c.removeIPToBl(entity.IPNetwork{
					IP:   setCommand[2],
					Mask: setCommand[3],
				})
			case "get":
				c.getIPListFromBl()
			default:
				fmt.Println("unknown command")
			}

		case "whitelist":
			switch setCommand[1] {
			case "add":
				if len(setCommand) != 4 {
					break
				}
				c.addIPToWl(entity.IPNetwork{
					IP:   setCommand[2],
					Mask: setCommand[3],
				})
			case "remove":
				if len(setCommand) != 4 {
					break
				}
				c.removeIPToWl(entity.IPNetwork{
					IP:   setCommand[2],
					Mask: setCommand[3],
				})
			case "get":
				c.getIPListFromWl()
			default:
				fmt.Println("unknown command")
			}

		case "bucket":
			if len(setCommand) != 4 {
				break
			}
			if setCommand[1] == "reset" {
				c.resetBucket(entity.Request{
					Login:    setCommand[2],
					Password: "",
					IP:       setCommand[3],
				})
			} else {
				fmt.Println("unknown command")
			}
		case "exit":
			ch <- os.Interrupt
			return
		case "help":
			for _, suggestion := range suggestions {
				fmt.Println("Command:", suggestion.Text, "Description:", suggestion.Description)
			}

		default:
			fmt.Println("unknown command")
		}
	})
	completer := prompt.Completer(func(in prompt.Document) []prompt.Suggest {
		w := in.GetWordBeforeCursor()
		if w == "" {
			return []prompt.Suggest{}
		}
		return prompt.FilterHasPrefix(suggestions, w, true)
	})
	defer func() {
		if a := recover(); a != nil {
			log.Println("Command line interface not available. Please run container with tty mode")
		}
	}()
	prompt.New(executor, completer).Run()
}

// addIPToBl добавляет IP-адрес в черный список.
func (c *CommandLineInterface) addIPToBl(ipNet entity.IPNetwork) {
	isValidateIP := handlers.ValidateIP(ipNet)
	if !isValidateIP {
		fmt.Println("not valid ip")
		return
	}
	err := c.serviceBL.AddIP(ipNet)
	if err != nil {
		fmt.Printf("service error: %v \n", err)
		return
	}
	fmt.Printf("add address: %v to blacklist \n", ipNet)
}

// removeIPToBl удаляет IP-адрес из черного списка.
func (c *CommandLineInterface) removeIPToBl(ipNet entity.IPNetwork) {
	isValidateIP := handlers.ValidateIP(ipNet)
	if !isValidateIP {
		fmt.Println("not valid ip")
		return
	}
	err := c.serviceBL.RemoveIP(ipNet)
	if err != nil {
		fmt.Printf("service error: %v \n", err)
		return
	}
	fmt.Printf("remove address: %v from blacklist \n", ipNet)
}

// getIPListFromBl получает список IP-адресов из черного списка и выводит их.
func (c *CommandLineInterface) getIPListFromBl() {
	list, err := c.serviceBL.GetIPList()
	if err != nil {
		return
	}
	for _, network := range list {
		fmt.Println(network)
	}
}

// addIPToWl добавляет IP-адрес в белый список.
func (c *CommandLineInterface) addIPToWl(ipNet entity.IPNetwork) {
	isValidateIP := handlers.ValidateIP(ipNet)
	if !isValidateIP {
		fmt.Println("not valid ip")
		return
	}
	err := c.serviceWL.AddIP(ipNet)
	if err != nil {
		fmt.Printf("service error: %v \n", err)
		return
	}
	fmt.Printf("add address: %v to whitelist \n", ipNet)
}

// removeIPToWl удаляет IP-адрес из белого списка.
func (c *CommandLineInterface) removeIPToWl(ipNet entity.IPNetwork) {
	isValidateIP := handlers.ValidateIP(ipNet)
	if !isValidateIP {
		fmt.Println("not valid ip")
		return
	}
	err := c.serviceWL.RemoveIP(ipNet)
	if err != nil {
		fmt.Printf("service error: %v \n", err)
		return
	}
	fmt.Printf("remove address: %v from whitelist \n", ipNet)
}

// getIPListFromWl получает список IP-адресов из белого списка и выводит их.
func (c *CommandLineInterface) getIPListFromWl() {
	list, err := c.serviceWL.GetIPList()
	if err != nil {
		return
	}
	for _, network := range list {
		fmt.Println(network)
	}
}

// resetBucket сбрасывает IP-адрес и/или логин из "бакета" (bucket).
func (c *CommandLineInterface) resetBucket(request entity.Request) {
	isValidateReq := handlers.ValidateRequest(request)
	if !isValidateReq {
		fmt.Println("not valid ip")
		return
	}
	isResetIP := c.serviceAuth.ResetIPBucket(request.IP)
	if !isResetIP {
		fmt.Printf("ip address: %v not find\n", request.IP)
	} else {
		fmt.Printf("ip address: %v has been reseted\n", request.IP)
	}

	isResetLogin := c.serviceAuth.ResetLoginBucket(request.IP)
	if !isResetLogin {
		fmt.Printf("login: %v not find\n", request.IP)
	} else {
		fmt.Printf("ip address: %v has been reseted\n", request.IP)
	}
}
