package agents

import (
	"net"
	"strings"
	"testing"
	"time"

	"github.com/cgrates/go-diameter/diam"
	"github.com/cgrates/go-diameter/diam/avp"
	"github.com/cgrates/go-diameter/diam/datatype"
	"github.com/cgrates/go-diameter/diam/dict"
	"github.com/cgrates/go-diameter/diam/sm"
)

func TestDiameterClient(t *testing.T) {
	// Setup a mock server on an ephemeral port
	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("Error listening: %v", err)
	}
	defer l.Close()
	addr := l.Addr().String()
	network := "tcp"

	srvMux := sm.New(&sm.Settings{
		OriginHost:       datatype.DiameterIdentity("server"),
		OriginRealm:      datatype.DiameterIdentity("server"),
		VendorID:         datatype.Unsigned32(10415),
		ProductName:      datatype.UTF8String("go-diameter"),
		FirmwareRevision: datatype.Unsigned32(1),
	})

	srv := &diam.Server{
		Network: network,
		Addr:    addr,
		Handler: srvMux,
	}

	// Note: CER is handled automatically by sm
	srvMux.HandleFunc("ALL", func(c diam.Conn, m *diam.Message) {
		// Any request other than CER
		if m.Header.CommandFlags&diam.RequestFlag != 0 {
			a := m.Answer(2001)
			a.NewAVP(avp.SessionID, avp.Mbit, 0, datatype.UTF8String("12345"))
			a.NewAVP(avp.OriginHost, avp.Mbit, 0, datatype.DiameterIdentity("server"))
			a.NewAVP(avp.OriginRealm, avp.Mbit, 0, datatype.DiameterIdentity("server"))
			_, err := a.WriteTo(c)
			if err != nil {
				t.Logf("Server failed to write answer: %v", err)
			}
		}
	})

	go func() {
		// diam.Server.Serve(net.Listener) blocks, it serves connections on the listener
		err := srv.Serve(l)
		if err != nil && !strings.Contains(err.Error(), "use of closed network connection") {
			t.Logf("Server exited: %v", err)
		}
	}()

	// Wait briefly to ensure server loop is ready
	time.Sleep(100 * time.Millisecond)

	client, err := NewDiameterClient(addr, "test_host", "test_realm", 10415, "test_product", 1, "", network)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}
	defer client.Close()

	if client.conn == nil {
		t.Errorf("Client connection should not be nil")
	}

	// Send a CreditControl request (CCR) to trigger the ALL handler
	m := diam.NewMessage(diam.CreditControl, diam.RequestFlag, 4, 0, 0, dict.Default)
	m.NewAVP(avp.SessionID, avp.Mbit, 0, datatype.UTF8String("12345"))
	m.NewAVP(avp.OriginHost, avp.Mbit, 0, datatype.DiameterIdentity("test_host"))
	m.NewAVP(avp.OriginRealm, avp.Mbit, 0, datatype.DiameterIdentity("test_realm"))
	m.NewAVP(avp.DestinationRealm, avp.Mbit, 0, datatype.DiameterIdentity("server"))

	err = client.SendMessage(m)
	if err != nil {
		t.Fatalf("Failed to send message: %v", err)
	}

	// Wait for message
	rcv := client.ReceivedMessage(2 * time.Second)
	if rcv == nil {
		t.Fatalf("Expected message, got nil")
	}

	// Try timeout
	rcv2 := client.ReceivedMessage(50 * time.Millisecond)
	if rcv2 != nil {
		t.Errorf("Expected nil on timeout, got %v", rcv2)
	}
}
