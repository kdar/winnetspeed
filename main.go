package main

import (
  "fmt"
  "github.com/StackExchange/wmi"
  "github.com/dustin/go-humanize"
  "github.com/nsf/termbox-go"
  "sync"
  "time"
)

type Win32_PerfFormattedData_Tcpip_NetworkInterface struct {
  Name                string
  BytesReceivedPerSec uint32
  BytesSentPerSec     uint32
}

func print(x, y int, s string) {
  for _, r := range s {
    termbox.SetCell(x, y, r, termbox.ColorDefault, termbox.ColorDefault)
    x += 1
  }
}

func main() {
  err := termbox.Init()
  if err != nil {
    panic(err)
  }
  defer termbox.Close()

  event_queue := make(chan termbox.Event)
  go func() {
    for {
      event_queue <- termbox.PollEvent()
    }
  }()

  var dst []Win32_PerfFormattedData_Tcpip_NetworkInterface
  q := wmi.CreateQuery(&dst, `` /*`WHERE Name = "Realtek PCIe GBE Family Controller"`*/)

  var wg sync.WaitGroup
  wg.Add(1)
  go func() {
    defer wg.Done()

  loop:
    for {
      select {
      case ev := <-event_queue:
        if ev.Type == termbox.EventKey && ev.Key == termbox.KeyEsc {
          break loop
        }
      default:
        var d []Win32_PerfFormattedData_Tcpip_NetworkInterface
        err := wmi.Query(q, &d)
        if err != nil {
          panic(err)
        }

        termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
        y := 0
        for _, entry := range d {
          print(0, y, fmt.Sprintf("%s: recv: %s sent: %s", entry.Name, humanize.Bytes(uint64(entry.BytesReceivedPerSec)), humanize.Bytes(uint64(entry.BytesSentPerSec))))
          y++
        }
        termbox.Flush()

        <-time.After(time.Millisecond * 500)
      }
    }
  }()

  wg.Wait()
}
