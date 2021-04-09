package main

import (
    "fmt"
    "os"

    "github.com/joho/godotenv"
    "github.com/hashgraph/hedera-sdk-go/v2"
)

type Account struct {
    private string
    publickey string
    accountid string
}

func CreateAccount(cli *hedera.Client) (*Account,error) {
    privk, err := hedera.GeneratePrivateKey()
    if err != nil {
        return nil, err
    }
    pubkey := privk.PublicKey()
    newAccountTx, err := hedera.NewAccountCreateTransaction().
        SetKey(pubkey).
        SetInitialBalance(hedera.HbarFrom(10, hedera.HbarUnits.Tinybar)).
        Execute(cli)
    if err != nil {
        return nil, err
    }
    receipt, err := newAccountTx.GetReceipt(cli)
    if err != nil {
        return nil, err
    }
    newAccountID := *receipt.AccountID
    return &Account{
        private: privk.String(),
        publickey: pubkey.String(),
        accountid: newAccountID.String(),
    }, nil

}

func main() {
    //Loads the .env file and throws an error if it cannot load the variables from that file corectly
    err := godotenv.Load(".env")
    if err != nil {
        panic(fmt.Errorf("Unable to load enviroment variables from .env file. Error:\n%v\n", err))
    }

    //Grab your testnet account ID and private key from the .env file
    myAccountId, err := hedera.AccountIDFromString(os.Getenv("MY_ACCOUNT_ID"))
    if err != nil {
        panic(err)
    }

    myPrivateKey, err := hedera.PrivateKeyFromString(os.Getenv("MY_PRIVATE_KEY"))
    if err != nil {
        panic(err)
    }

    //Print your testnet account ID and private key to the console to make sure there was no error
    fmt.Printf("The account ID is = %v\n", myAccountId)
    fmt.Printf("The private key is = %v\n", myPrivateKey)
    client := hedera.ClientForTestnet()
    client.SetOperator(myAccountId,myPrivateKey)

    query := hedera.NewAccountBalanceQuery().
        SetAccountID(myAccountId)
    myBalance, err := query.Execute(client)
    if err != nil {
        panic(err)
    }
    fmt.Println("The account balance is ", myBalance.Hbars.AsTinybar())

    //account, err := CreateAccount(client)
    //if err != nil {
    //    panic(err)
    //}
    //fmt.Printf("New account %v\n", account)

}
