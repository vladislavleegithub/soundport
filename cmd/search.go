package cmd

import (
	"fmt"

	"github.com/Samarthbhat52/soundport/api/ytmusic"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(searchCmd)
}

var searchCmd = &cobra.Command{
	Use: "search",
	Run: func(cmd *cobra.Command, args []string) {
		val, err := ytmusic.SearchSongYT(
			[]string{
				"Blinding Lights, The Weeknd",
				"Levitating, Dua Lipa",
				"Good 4 U, Olivia Rodrigo",
				"Shape of You, Ed Sheeran",
				"Stay, Kid LAROI & Justin Bieber",
				"Watermelon Sugar, Harry Styles",
				"Bad Habits, Ed Sheeran",
				"Save Your Tears, The Weeknd",
				"Montero (Call Me By Your Name), Lil Nas X",
				"Peaches, Justin Bieber ft. Daniel Caesar, Giveon",
				"Kiss Me More, Doja Cat ft. SZA",
				"Heat Waves, Glass Animals",
				"Industry Baby, Lil Nas X & Jack Harlow",
				"Shivers, Ed Sheeran",
				"As It Was, Harry Styles",
				"Easy On Me, Adele",
				"Deja Vu, Olivia Rodrigo",
				"Astronaut in the Ocean, Masked Wolf",
				"Good Days, SZA",
				"Take My Breath, The Weeknd",
				"Blame It On Me, George Ezra",
				"Rolling in the Deep, Adele",
				"Titanium, David Guetta ft. Sia",
				"Uptown Funk, Mark Ronson ft. Bruno Mars",
				"Someone Like You, Adele",
				"Sunflower, Post Malone & Swae Lee",
				"All of Me, John Legend",
				"Don’t Start Now, Dua Lipa",
				"Shallow, Lady Gaga & Bradley Cooper",
				"Old Town Road, Lil Nas X ft. Billy Ray Cyrus",
				"Happier, Marshmello ft. Bastille",
				"We Found Love, Rihanna ft. Calvin Harris",
				"Faded, Alan Walker",
				"Perfect, Ed Sheeran",
				"Rockstar, Post Malone ft. 21 Savage",
				"Call Me Maybe, Carly Rae Jepsen",
				"Blow, Ed Sheeran, Bruno Mars & Chris Stapleton",
				"In the Name of Love, Martin Garrix & Bebe Rexha",
				"Savage Love, Jawsh 685, Jason Derulo",
				"Closer, The Chainsmokers ft. Halsey",
				"The Middle, Zedd, Maren Morris, Grey",
				"Lose Yourself, Eminem",
				"Starboy, The Weeknd ft. Daft Punk",
				"Don’t Let Me Down, The Chainsmokers ft. Daya",
				"Take A Bow, Rihanna",
				"Rolling Stone, The Rolling Stones",
				"Toxic, Britney Spears",
				"Viva La Vida, Coldplay",
				"Counting Stars, OneRepublic",
			},
		)
		if err != nil {
			return
		}

		fmt.Println(val)
	},
}
