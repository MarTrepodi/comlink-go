package main

import (
	"fmt"
	"log"
	"swgoh-comlink-client/swgohcomlink"
)

func main() {
	// Initialize the client. Replace with the correct API base URL.
	client := swgohcomlink.NewClient(swgohcomlink.DefaultBaseURL)

	// --- Example 1: GetMetaData ---
	fmt.Println("--- Fetching Metadata ---")
	metaReq := &swgohcomlink.GetMetaDataRequest{
		Payload: &swgohcomlink.GetMetaDataPayload{
			ClientSpecs: &swgohcomlink.GetMetaDataClientSpecs{
				Platform: "Android",
			},
		},
	}
	metadata, err := client.GetMetaData(metaReq)
	if err != nil {
		// This will likely fail with the placeholder URL: log.Printf("Error fetching metadata: %v", err)
		log.Printf("Error (expected with placeholder URL) fetching metadata: %v", err)
	} else {
		fmt.Printf("Metadata Version: %s\n", metadata.LatestGamedataVersion)
		fmt.Printf("Server Timestamp: %d\n", metadata.ServerTimestamp)
	}

	fmt.Println("\n--- Example 2: GetPlayer ---")
	// Replace "123456789" with a valid ally code for a real test.
	playerReq := &swgohcomlink.GetPlayerRequest{
		Payload: swgohcomlink.GetPlayerPayload{
			AllyCode: "123456789",
		},
	}

	// Check for a valid request first
	if err := playerReq.Validate(); err != nil {
		log.Fatalf("Invalid player request: %v", err)
	}

	player, err := client.GetPlayer(playerReq)
	if err != nil {
		// This will likely fail with the placeholder URL: log.Printf("Error fetching player: %v", err)
		log.Printf("Error (expected with placeholder URL) fetching player: %v", err)
	} else {
		fmt.Printf("Player Name: %s\n", player.Name)
		fmt.Printf("Guild Name: %s\n", player.GuildName)
		fmt.Printf("Total Units: %d\n", len(player.RosterUnit))
	}
}
