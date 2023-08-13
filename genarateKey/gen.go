package main

import (
	"encoding/base64"
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
)

func main() {
	pwd, err := os.Getwd()
	if err != nil {
		logrus.Error("can't get working directory err: ", err)

	}

	privateKey, err := os.ReadFile(pwd + "/cert/id_rsa")
	if err != nil {
		logrus.Error("can't read private key err: ", err)

	}

	publicKey, err := os.ReadFile(pwd + "/cert/id_rsa.pub")
	if err != nil {
		logrus.Error("can't read public key err: ", err)

	}

	f, err := os.Create(".env")
	if err != nil {
		logrus.Error("can't create file err: ", err)
	}
	defer f.Close()

	// return privateKey, publicKey, nil
	private := base64.StdEncoding.EncodeToString(privateKey)
	public := base64.StdEncoding.EncodeToString(publicKey)

	// fmt.Println("privateKey : ", private)
	// fmt.Println("publicKey : ", public)

	_, err = f.WriteString("PUBLIC_KEY=" + public + "\n")
	if err != nil {
		logrus.Error("can't write to file err: ", err)
	}

	_, err = f.WriteString("PRIVATE_KEY=" + private + "\n")
	if err != nil {
		logrus.Error("can't write to file err: ", err)
	}

	mongoUri := "MONGO_URL=mongodb://localhost:27017,localhost:27018,localhost:27019/my-database?replicaSet=my-replica-set"
	_, err = f.WriteString(mongoUri + "\n")
	if err != nil {
		logrus.Error("can't write to file err: ", err)
	}
	fmt.Println("done")

}
