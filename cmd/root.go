/*
Copyright © 2019 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
  "fmt"
  "github.com/mitchellh/go-homedir"
  "github.com/spf13/cobra"
  "github.com/spf13/viper"
  "os"
)


var cfgFile string
var token string
var username string
var password string
var hostname string
var verbose bool
var allowedUsers []string
var https bool
var port int

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
  Use:   "transmission-telegram-bot",
  Short: "A Telegram chatbot for managing torrents through Transmission",
  Long: `A Telegram chatbot that provides basic torrent management in Transmission.
It can add, stop, resume, and delete torrents, and supports limiting access to only
certain users.`,
  // Uncomment the following line if your bare application
  // has an action associated with it:
  //	Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
  if err := rootCmd.Execute(); err != nil {
    fmt.Println(err)
    os.Exit(1)
  }
}

func init() {
  cobra.OnInitialize(initConfig)

  // Here you will define your flags and configuration settings.
  // Cobra supports persistent flags, which, if defined here,
  // will be global for your application.

  rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.transmission-telegram-bot.json)")
  rootCmd.PersistentFlags().StringVar(&token, "bot-token", "", "authentication token for Telegram bot API")
  rootCmd.PersistentFlags().StringVar(&username, "username", "", "username for Transmission remote server (only use if auth is enabled)")
  rootCmd.PersistentFlags().StringVar(&password, "password", "", "password for Transmission remote server (only use if auth is enabled)")
  rootCmd.PersistentFlags().StringVar(&hostname, "hostname", "", "hostname for Transmission remote server")
  rootCmd.PersistentFlags().StringSliceVar(&allowedUsers, "allowed-user", []string{}, "Telegram username of user authorized to access this bot instance (defaults to no one)")
  rootCmd.PersistentFlags().BoolVar(&verbose, "verbose", false, "run bot with verbose log messages")
  rootCmd.PersistentFlags().BoolVar(&https, "https", true, "use HTTPS for Transmission RPC calls (defaults to true)")
  rootCmd.PersistentFlags().IntVar(&port, "port", 9091, "port for Transmission remote server (default is 9091)")

  // Cobra also supports local flags, which will only run
  // when this action is called directly.
  rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}


// initConfig reads in config file and ENV variables if set.
func initConfig() {
  if cfgFile != "" {
    // Use config file from the flag.
    viper.SetConfigFile(cfgFile)
  } else {
    // Find home directory.
    home, err := homedir.Dir()
    if err != nil {
      fmt.Println(err)
      os.Exit(1)
    }

    // Search config in home directory with name ".app" (without extension).
    viper.AddConfigPath(home)
    viper.SetConfigName(".transmission-telegram-bot")
    viper.SetConfigType("json")
  }

  _ = viper.BindPFlag("bot-token", rootCmd.PersistentFlags().Lookup("bot-token"))
  _ = viper.BindPFlag("username", rootCmd.PersistentFlags().Lookup("username"))
  _ = viper.BindPFlag("password", rootCmd.PersistentFlags().Lookup("password"))
  _ = viper.BindPFlag("hostname", rootCmd.PersistentFlags().Lookup("hostname"))
  _ = viper.BindPFlag("allowed-users", rootCmd.PersistentFlags().Lookup("allowed-user"))
  _ = viper.BindPFlag("verbose", rootCmd.PersistentFlags().Lookup("verbose"))
  _ = viper.BindPFlag("https", rootCmd.PersistentFlags().Lookup("https"))
  _ = viper.BindPFlag("port", rootCmd.PersistentFlags().Lookup("port"))

  // If a config file is found, read it in.
  if err := viper.ReadInConfig(); err == nil {
    fmt.Println("Using config file:", viper.ConfigFileUsed())
  }
}

