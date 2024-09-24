package stdout

import (
	"fmt"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
)

type stdout struct {
}

func New() (stdout, error) {
	return stdout{}, nil
}

func (s stdout) MyName() string {
	return "stdout"
}

func (s stdout) Write(packet gopacket.Packet, dnsLayer *layers.DNS) error {
	fmt.Println("----------------------")
	fmt.Println(packet.Metadata().Timestamp)
	fmt.Println("questions")
	for i, q := range dnsLayer.Questions {
		fmt.Println(i+1, string(q.Name), q.Type.String(), q.Class.String())
	}
	if len(dnsLayer.Answers) > 0 {
		fmt.Println("answers")
		for i, a := range dnsLayer.Answers {
			fmt.Println(i+1, string(a.Name), a.IP.String(), a.Type.String(), a.Class.String())
		}
	}
	fmt.Println()
	return nil
}

func (s stdout) Close() error {
	return nil
}
