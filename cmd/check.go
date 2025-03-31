package cmd

import (
	"fmt"
	"time"

	"github.com/aethiopicuschan/memcached-checker/internal/logger"
	"github.com/bradfitz/gomemcache/memcache"
	"github.com/briandowns/spinner"

	"github.com/spf13/cobra"
)

var checkCmd = &cobra.Command{
	Use:   "check",
	Short: "Check that it operates correctly as both a memcached and a compatible server.",
	Long:  "Check that it operates correctly as both a memcached and a compatible server.",
	RunE:  check,
}

func init() {
	checkCmd.Flags().StringP("address", "a", "127.0.0.1:11211", "Address of the memcached server")
	checkCmd.Flags().BoolP("flush", "f", false, "Flush the memcached server before running the check")

	rootCmd.AddCommand(checkCmd)
}

func check(cmd *cobra.Command, args []string) (err error) {
	mc := memcache.New(cmd.Flag("address").Value.String())
	defer mc.Close()

	// Ping
	s := spinner.New(spinner.CharSets[14], 100*time.Millisecond, spinner.WithSuffix(" Ping"))
	s.Start()
	if err = mc.Ping(); err != nil {
		logger.Fail("Ping")
		return
	}
	s.Stop()
	logger.OK("Ping")

	// Flush
	if cmd.Flag("flush").Changed {
		s = spinner.New(spinner.CharSets[14], 100*time.Millisecond, spinner.WithSuffix(" Flush"))
		s.Start()
		if err = mc.FlushAll(); err != nil {
			logger.Fail("Flush")
			return
		}
		s.Stop()
		logger.OK("Flush")
	}

	// Set
	s = spinner.New(spinner.CharSets[14], 100*time.Millisecond, spinner.WithSuffix(" Set"))
	s.Start()
	if err = mc.Set(&memcache.Item{Key: "key_set", Value: []byte("value_set")}); err != nil {
		s.Stop()
		logger.Fail("Set")
		return
	}
	item, err := mc.Get("key_set")
	if err != nil {
		s.Stop()
		logger.Fail("Get")
		return
	}
	if string(item.Value) != "value_set" {
		s.Stop()
		logger.Fail("Get invalid")
		err = fmt.Errorf("expected `value_set`, got `%s`", string(item.Value))
		return
	}
	s.Stop()
	logger.OK("Set")

	// Add
	s = spinner.New(spinner.CharSets[14], 100*time.Millisecond, spinner.WithSuffix(" Add"))
	s.Start()
	if err = mc.Add(&memcache.Item{Key: "key_add", Value: []byte("value_add")}); err != nil {
		s.Stop()
		logger.Fail("Add")
		return
	}
	item, err = mc.Get("key_add")
	if err != nil {
		s.Stop()
		logger.Fail("Get")
		return
	}
	if string(item.Value) != "value_add" {
		s.Stop()
		logger.Fail("Get invalid")
		err = fmt.Errorf("expected `value_add`, got `%s`", string(item.Value))
		return
	}
	s.Stop()
	logger.OK("Add")

	// Replace
	s = spinner.New(spinner.CharSets[14], 100*time.Millisecond, spinner.WithSuffix(" Replace"))
	s.Start()
	if err = mc.Replace(&memcache.Item{Key: "key_set", Value: []byte("value_replaced")}); err != nil {
		s.Stop()
		logger.Fail("Replace")
		return
	}
	item, err = mc.Get("key_set")
	if err != nil {
		s.Stop()
		logger.Fail("Get")
		return
	}
	if string(item.Value) != "value_replaced" {
		s.Stop()
		logger.Fail("Get invalid")
		err = fmt.Errorf("expected `value_replaced`, got `%s`", string(item.Value))
		return
	}
	s.Stop()
	logger.OK("Replace")

	// Get
	s = spinner.New(spinner.CharSets[14], 100*time.Millisecond, spinner.WithSuffix(" Get"))
	s.Start()
	item, err = mc.Get("key_set")
	if err != nil {
		s.Stop()
		logger.Fail("Get")
		return
	}
	if string(item.Value) != "value_replaced" {
		s.Stop()
		logger.Fail("Get invalid")
		err = fmt.Errorf("expected `value_replaced`, got `%s`", string(item.Value))
		return
	}
	s.Stop()
	logger.OK("Get")

	// Gets
	s = spinner.New(spinner.CharSets[14], 100*time.Millisecond, spinner.WithSuffix(" Gets"))
	s.Start()
	items, err := mc.GetMulti([]string{"key_set", "key_add"})
	if err != nil {
		s.Stop()
		return
	}
	if len(items) != 2 {
		s.Stop()
		logger.Fail("Gets length")
		err = fmt.Errorf("expected 2 items, got %d", len(items))
		return
	}
	if string(items["key_set"].Value) != "value_replaced" {
		s.Stop()
		logger.Fail("Gets invalid")
		err = fmt.Errorf("expected `value_replaced`, got `%s`", string(items["key_set"].Value))
		return
	}
	if string(items["key_add"].Value) != "value_add" {
		s.Stop()
		logger.Fail("Gets invalid")
		err = fmt.Errorf("expected `value_add`, got `%s`", string(items["key_add"].Value))
		return
	}
	s.Stop()
	logger.OK("Gets")

	// Append
	s = spinner.New(spinner.CharSets[14], 100*time.Millisecond, spinner.WithSuffix(" Append"))
	s.Start()
	if err = mc.Append(&memcache.Item{Key: "key_set", Value: []byte("_appended")}); err != nil {
		s.Stop()
		logger.Fail("Append")
		return
	}
	item, err = mc.Get("key_set")
	if err != nil {
		s.Stop()
		logger.Fail("Get")
		return
	}
	if string(item.Value) != "value_replaced_appended" {
		s.Stop()
		logger.Fail("Get invalid")
		err = fmt.Errorf("expected `value_replaced_appended`, got `%s`", string(item.Value))
		return
	}
	s.Stop()
	logger.OK("Append")

	// Prepend
	s = spinner.New(spinner.CharSets[14], 100*time.Millisecond, spinner.WithSuffix(" Prepend"))
	s.Start()
	if err = mc.Prepend(&memcache.Item{Key: "key_set", Value: []byte("prepended_")}); err != nil {
		s.Stop()
		logger.Fail("Prepend")
		return
	}
	item, err = mc.Get("key_set")
	if err != nil {
		s.Stop()
		logger.Fail("Get")
		return
	}
	if string(item.Value) != "prepended_value_replaced_appended" {
		s.Stop()
		logger.Fail("Get invalid")
		err = fmt.Errorf("expected `prepended_value_replaced_appended`, got `%s`", string(item.Value))
		return
	}
	s.Stop()
	logger.OK("Prepend")

	// Number(Incr, Decr)
	s = spinner.New(spinner.CharSets[14], 100*time.Millisecond, spinner.WithSuffix(" Increment"))
	s.Start()
	if err = mc.Set(&memcache.Item{Key: "key_number", Value: []byte("1")}); err != nil {
		s.Stop()
		logger.Fail("Set")
		return
	}
	if _, err = mc.Increment("key_number", 1); err != nil {
		s.Stop()
		logger.Fail("Increment")
		return
	}
	item, err = mc.Get("key_number")
	if err != nil {
		s.Stop()
		logger.Fail("Get")
		return
	}
	if string(item.Value) != "2" {
		s.Stop()
		logger.Fail("Get invalid")
		err = fmt.Errorf("expected `2`, got `%s`", string(item.Value))
		return
	}
	s.Stop()
	logger.OK("Increment")
	s = spinner.New(spinner.CharSets[14], 100*time.Millisecond, spinner.WithSuffix(" Decrement"))
	s.Start()
	if _, err = mc.Decrement("key_number", 1); err != nil {
		s.Stop()
		logger.Fail("Decrement")
		return
	}
	item, err = mc.Get("key_number")
	if err != nil {
		s.Stop()
		logger.Fail("Get")
		return
	}
	if string(item.Value) != "1" {
		s.Stop()
		logger.Fail("Get invalid")
		err = fmt.Errorf("expected `1`, got `%s`", string(item.Value))
		return
	}
	s.Stop()
	logger.OK("Decrement")

	// Touch
	s = spinner.New(spinner.CharSets[14], 100*time.Millisecond, spinner.WithSuffix(" Touch"))
	s.Start()
	if err = mc.Set(&memcache.Item{Key: "key_touch", Value: []byte("value_touch")}); err != nil {
		logger.Fail("Set")
		return
	}
	if err = mc.Touch("key_touch", 2); err != nil {
		logger.Fail("Touch")
		return
	}
	time.Sleep(3 * time.Second)
	s.Stop()
	_, err = mc.Get("key_touch")
	if err == memcache.ErrCacheMiss {
		logger.OK("Touch")
		err = nil
	} else {
		logger.Fail("Touch")
		err = fmt.Errorf("expected `ErrNotStored`, got `%v`", err)
		return
	}

	// Flush
	s = spinner.New(spinner.CharSets[14], 100*time.Millisecond, spinner.WithSuffix(" Flush"))
	s.Start()
	if err = mc.FlushAll(); err != nil {
		s.Stop()
		logger.Fail("Flush")
		return
	}
	items, err = mc.GetMulti([]string{"key_set", "key_add"})
	s.Stop()
	if len(items) != 0 {
		logger.Fail("Flush")
		err = fmt.Errorf("expected empty")
		return
	}
	logger.OK("Flush")

	return
}
