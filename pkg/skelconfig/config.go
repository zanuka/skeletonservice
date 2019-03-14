package skelconfig

import (
	"bytes"
	"cloud.google.com/go/kms/apiv1"
	"cloud.google.com/go/storage"
	"context"
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/spf13/viper"
	kmspb "google.golang.org/genproto/googleapis/cloud/kms/v1"
	"io/ioutil"
	"log"
)

var AppConfig appConfig

var DBconn *sql.DB

type appConfig struct {
	DBHost     string
	DBUser     string
	DBPassword string
	DBName     string
	DBPort     int

	HashPepper string
}

var kmsClient *kms.KeyManagementClient

// init is auto run on import - this sets up the whole config for the app
// It also check for all the runtime vars and db connections required
// to make the auth service run
func init() {

	viper.SetEnvPrefix("skelsvc")
	viper.AutomaticEnv()

	bucket := viper.Get("CONFIG_BUCKET")
	configLocation := viper.Get("CONFIG_LOCATION")

	var decodedConf []byte
	var err error
	if bucket == nil || configLocation == nil {
		log.Println("Env vars SKELSVC_CONFIG_BUCKET and/or SKELSVC_CONFIG_LOCATION not found")
		decodedConf = loadConfigFromLocal()
	} else {
		decodedConf, err = loadConfigFromStorage(bucket.(string), configLocation.(string))
		if err != nil {
			log.Fatal(err)
		}
	}

	// load the config into viper
	viper.SetConfigType("yaml")
	err = viper.ReadConfig(bytes.NewReader(decodedConf))
	if err != nil { // Handle errors reading the config file
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}

	// make is easily easy to access and check all the vars
	checkAndParseConfig()
	checkDBConnection()
}

func checkDBConnection() {

	dbinfo := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d",
		AppConfig.DBHost,
		AppConfig.DBUser,
		AppConfig.DBPassword,
		AppConfig.DBName,
		AppConfig.DBPort)

	var err error
	DBconn, err = sql.Open("postgres", dbinfo)
	if err != nil {
		log.Fatal(err)
	}

	err = DBconn.Ping()
	if err != nil {
		DBconn.Close()
		log.Fatal(err)
	}

}

func loadConfigFromLocal() []byte {

	log.Println("Attempting to load config from local file: local-config.yaml.enc ")

	dat, err := ioutil.ReadFile("local-config.yaml.enc")
	if err != nil {
		log.Fatal(err)
	}

	decodedConf, err1 := decryptConfig(dat)
	if err1 != nil {
		log.Fatal(err1)
		//return nil, err
	}

	return decodedConf

}

// explicit reads credentials from the specified path.
func loadConfigFromStorage(bucket string, configPath string) ([]byte, error) {

	log.Println(fmt.Sprintf("Loading config from cloud storage. path=%v/%v", bucket, configPath))

	ctx := context.Background()

	// For API packages whose import path is starting with "cloud.google.com/go",
	// such as cloud.google.com/go/storage in this case, if there are no credentials
	// provided, the client library will look for credentials in the environment.
	storageClient, err := storage.NewClient(ctx)
	if err != nil {
		log.Fatal(err)
	}

	rc, err := storageClient.Bucket(bucket).Object(configPath).NewReader(ctx)
	if err != nil {
		log.Fatal(err)
		//return nil, err
	}
	defer rc.Close()

	slurp, err := ioutil.ReadAll(rc)
	if err != nil {
		log.Fatal(err)
	}

	decodedConf, err1 := decryptConfig(slurp)
	if err1 != nil {
		log.Fatal(err1)
	}

	return decodedConf, nil

}

func checkAndParseConfig() {

	AppConfig = appConfig{}

	// DB VARS
	AppConfig.DBHost = viper.GetString("DB_HOST")
	AppConfig.DBUser = viper.GetString("DB_USER")
	AppConfig.DBPassword = viper.GetString("DB_PASSWORD")
	AppConfig.DBName = viper.GetString("DB_NAME")
	AppConfig.DBPort = viper.GetInt("DB_PORT")

	// HASH VARS
	AppConfig.HashPepper = viper.GetString("HASH_PEPPER")

	if AppConfig.DBHost == "" {
		log.Fatal("Cannot find the DB_HOST in the config")
	}

	if AppConfig.DBUser == "" {
		log.Fatal("Cannot find DB_USER in the config")
	}

	if AppConfig.DBPassword == "" {
		log.Fatal("Cannot find DB_PASSWORD in the config")
	}

	if AppConfig.DBName == "" {
		log.Fatal("Cannot find DB_NAME in the config")
	}

	if AppConfig.DBName == "" {
		log.Fatal("Cannot find DB_NAME in the config")
	}

	if AppConfig.DBPort == 0 {
		log.Fatal("Cannot find DB_PORT in the config")
	}

	if AppConfig.HashPepper == "" {
		log.Fatal("Cannot find HASH_PEPPER in the config")
	}
}

func decryptConfig(cipher []byte) ([]byte, error) {
	ctx := context.Background()

	var err error
	kmsClient, err = kms.NewKeyManagementClient(ctx)
	if err != nil {

		s := `
-----------------------------------------------------------------------------
APPLICATION START ERROR
------------------------------------------------------------------------------
DOING LOCAL DEV?:
Ask for a service account or from the Cloud Console, create a service account, 
download its json credentials file, then set the GOOGLE_APPLICATION_CREDENTIALS 
environment variable:
export GOOGLE_APPLICATION_CREDENTIALS=/path/to/your-project-credentials.json
-----------------------------------------------------------------------------
`
		log.Fatal(s)
	}

	req := &kmspb.DecryptRequest{
		Name:       "projects/[project-id]/locations/global/keyRings/[keyring]/cryptoKeys/[key]",
		Ciphertext: cipher,
	}

	resp, err := kmsClient.Decrypt(ctx, req)

	if err != nil {
		return nil, err
	}
	return resp.Plaintext, nil

}
