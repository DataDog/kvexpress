// +build linux darwin freebsd

package commands

import (
	"fmt"
	consul "github.com/hashicorp/consul/api"
	"github.com/spf13/cobra"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var (
	// WatchPrefix is the Consul KV space to watch.
	WatchPrefix string

	// WriteLockKey is the key where kvexpress writes the lock data.
	// We are using PrefixLocation so that you can have multiple readers with
	// different read locks but only a single writer per KV destination.
	WriteLockKey string

	// LeaderCh tells us whether or not we've got the lock.
	LeaderCh <-chan struct{}

	// Lock is for a LockKey
	Lock *consul.Lock
)

const (
	watchSleep = 2
)

var watchCmd = &cobra.Command{
	Use:   "watch",
	Short: "Watch a prefix for changed keys.",
	Long:  `Watch a prefix for changed keys. Convert incoming keys to kvexpress format and location.`,
	PreRun: func(cmd *cobra.Command, args []string) {
		checkWatchFlags()
		AutoEnable()
	},
	Run: watchRun,
}

func watchRun(cmd *cobra.Command, args []string) {
	// Check and aquire lock.
	c, err := Connect(ConsulServer, Token)
	if err != nil {
		Log("Could not connect to consul.", "info")
	}
	Lock, err = c.LockKey(WriteLockKey)
	if err != nil {
		Log("Could not setup LockKey.", "info")
	}
	go acquireConsulLock(c)

	// Setup a goroutine to teardown the lock properly.
	ctrlc := make(chan os.Signal)
	signal.Notify(ctrlc, os.Interrupt, syscall.SIGTERM)
	go teardownLock(ctrlc)

	// When we've got a lock - then we can setup a watch.
	go startWatch()

	time.Sleep(time.Duration(300) * time.Second)
	// If you don't get a lock - then wait until you do.
	// If you've got a lock - set a watch on WatchPrefix.

	// If the watch fires - get the changed key and write into the PrefixLocation heirarchy.
}

func checkWatchFlags() {
	Log("Checking cli flags.", "debug")
	if WatchPrefix == "" {
		fmt.Println("Need a KV space to watch in -w")
		os.Exit(1)
	}
	if WatchPrefix == PrefixLocation {
		fmt.Println("You cannot watch the same location you're going to store the result in.")
		os.Exit(1)
	}
	if PrefixLocation != "" {
		WriteLockKey = fmt.Sprintf("%s/.lock", PrefixLocation)
	}
	Log("Required cli flags present.", "debug")
	fmt.Printf("Starting up - Control-C to exit.\n")
}

func init() {
	RootCmd.AddCommand(watchCmd)
	watchCmd.Flags().StringVarP(&WatchPrefix, "watch", "w", "", "Consul KV space to watch.")
}

func startWatch() {
	for {
		if LeaderCh == nil {
			Log("I do NOT have the lock - waiting.", "info")
		} else {
			Log("I have the lock - let's setup a watch.", "info")
		}
		time.Sleep(time.Duration(5) * time.Second)
	}
}

func acquireConsulLock(c *consul.Client) {
	var err error
	for {
		if LeaderCh == nil {
			Log("Trying to acquire a Consul lock.", "info")
			LeaderCh, err = Lock.Lock(nil)
			if LeaderCh != nil {
				Log("I have aquired a lock.", "info")
			}
			if err != nil {
				Log(fmt.Sprintf("err: %v", err), "info")
				Log("I do NOT have the lock.", "info")
				time.Sleep(time.Duration(watchSleep) * time.Second)
				acquireConsulLock(c)
			}
		} else {
			Log("Already have a lock - not reacquiring", "info")
		}
		time.Sleep(time.Duration(watchSleep) * time.Second)
	}
}

// Let's properly teardown the lock
func teardownLock(c chan os.Signal) {
	sig := <-c
	message := fmt.Sprintf("Received '%s' - shutting down.", sig)
	Log(message, "info")
	fmt.Printf("%s\n", message)

	// Unlock Consul.
	fmt.Printf("Unlocking Consul.\n")
	err := Lock.Unlock()
	if err != nil {
		Log(fmt.Sprintf("Could not unlock. Err: %s", err), "info")
		os.Exit(1)
	}
	fmt.Printf("Consul unlocked.\n")
	os.Exit(0)
}
