package main

import (
	"context"
	"errors"
	"fmt"
	"io"

	"github.com/AntonNikol/anti-bruteforce/internal/config"
	"github.com/AntonNikol/anti-bruteforce/internal/controller/grpcapi/authorizationpb"
	"github.com/AntonNikol/anti-bruteforce/internal/controller/grpcapi/blacklistpb"
	"github.com/AntonNikol/anti-bruteforce/internal/controller/grpcapi/bucketpb"
	"github.com/AntonNikol/anti-bruteforce/internal/controller/grpcapi/whitelistpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	c, err := config.LoadAll() // Инициализация конфигурации приложения.
	if err != nil {
		fmt.Println(err) // Вывод ошибки и завершение программы в случае неудачи.
		return
	}

	insecureTr := grpc.WithTransportCredentials(insecure.NewCredentials())
	dial, err := grpc.Dial(
		c.Listen.BindIP+":"+c.Listen.Port, insecureTr) // Установка соединения с удаленным gRPC-сервером.
	if err != nil {
		fmt.Println(err) // Вывод ошибки и завершение программы в случае неудачи.
		return
	}

	// Создание клиентов для различных gRPC-сервисов
	clientBL := blacklistpb.NewBlackListServiceClient(dial)
	clientWL := whitelistpb.NewWhiteListServiceClient(dial)
	clientBucket := bucketpb.NewBucketServiceClient(dial)
	clientAuth := authorizationpb.NewAuthorizationClient(dial)

	getIPListInBlackList(clientBL) // Вызов функции для получения списка IP-адресов из черного списка.
	fmt.Println()
	getIPListInWhiteList(clientWL) // Вызов функции для получения списка IP-адресов из белого списка.
	fmt.Println()
	resetBucket(clientBucket) // Вызов функции для сброса бакета.
	fmt.Println()
	tryAuth(clientAuth) // Вызов функции для попытки авторизации.
}

// Функция для попытки авторизации.
func tryAuth(client authorizationpb.AuthorizationClient) {
	response, err := client.TryAuthorization(
		context.Background(),
		&authorizationpb.AuthorizationRequest{Request: &authorizationpb.Request{
			Login:    "test",
			Password: "1234",
			Ip:       "192.1.5.1",
		}})
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(response.IsAllow)
}

// Функция для получения списка IP-адресов из черного списка.
func getIPListInBlackList(client blacklistpb.BlackListServiceClient) {
	stream, err := client.GetIpList(context.Background(), &blacklistpb.GetIpListRequest{})
	if err != nil {
		fmt.Println(err)
		return
	}
	for {
		res, err := stream.Recv()

		if errors.Is(err, io.EOF) {
			return
		}
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(res)
	}
}

// Функция для получения списка IP-адресов из белого списка.
func getIPListInWhiteList(client whitelistpb.WhiteListServiceClient) {
	stream, err := client.GetIpList(context.Background(), &whitelistpb.GetIpListRequest{})
	if err != nil {
		fmt.Println(err)
		return
	}
	for {
		res, err := stream.Recv()
		if errors.Is(err, io.EOF) {
			return
		}
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(res)
	}
}

// Функция для сброса бакета.
func resetBucket(client bucketpb.BucketServiceClient) {
	response, err := client.ResetBucket(context.Background(), &bucketpb.ResetBucketRequest{Request: &bucketpb.Request{
		Login:    "test",
		Password: "1234",
		Ip:       "192.1.5.1",
	}})
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(response.ResetIp, response.ResetLogin)
}
