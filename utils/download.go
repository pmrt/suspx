package utils

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
)

const MaxConcurrentDownloads = 3

var csvParts = []string{
	"https://placedata.reddit.com/data/canvas-history/2022_place_canvas_history-000000000000.csv.gzip",
	"https://placedata.reddit.com/data/canvas-history/2022_place_canvas_history-000000000001.csv.gzip",
	"https://placedata.reddit.com/data/canvas-history/2022_place_canvas_history-000000000002.csv.gzip",
	"https://placedata.reddit.com/data/canvas-history/2022_place_canvas_history-000000000003.csv.gzip",
	"https://placedata.reddit.com/data/canvas-history/2022_place_canvas_history-000000000004.csv.gzip",
	"https://placedata.reddit.com/data/canvas-history/2022_place_canvas_history-000000000005.csv.gzip",
	"https://placedata.reddit.com/data/canvas-history/2022_place_canvas_history-000000000006.csv.gzip",
	"https://placedata.reddit.com/data/canvas-history/2022_place_canvas_history-000000000007.csv.gzip",
	"https://placedata.reddit.com/data/canvas-history/2022_place_canvas_history-000000000008.csv.gzip",
	"https://placedata.reddit.com/data/canvas-history/2022_place_canvas_history-000000000009.csv.gzip",
	"https://placedata.reddit.com/data/canvas-history/2022_place_canvas_history-000000000010.csv.gzip",
	"https://placedata.reddit.com/data/canvas-history/2022_place_canvas_history-000000000011.csv.gzip",
	"https://placedata.reddit.com/data/canvas-history/2022_place_canvas_history-000000000012.csv.gzip",
	"https://placedata.reddit.com/data/canvas-history/2022_place_canvas_history-000000000013.csv.gzip",
	"https://placedata.reddit.com/data/canvas-history/2022_place_canvas_history-000000000014.csv.gzip",
	"https://placedata.reddit.com/data/canvas-history/2022_place_canvas_history-000000000015.csv.gzip",
	"https://placedata.reddit.com/data/canvas-history/2022_place_canvas_history-000000000016.csv.gzip",
	"https://placedata.reddit.com/data/canvas-history/2022_place_canvas_history-000000000017.csv.gzip",
	"https://placedata.reddit.com/data/canvas-history/2022_place_canvas_history-000000000018.csv.gzip",
	"https://placedata.reddit.com/data/canvas-history/2022_place_canvas_history-000000000019.csv.gzip",
	"https://placedata.reddit.com/data/canvas-history/2022_place_canvas_history-000000000020.csv.gzip",
	"https://placedata.reddit.com/data/canvas-history/2022_place_canvas_history-000000000021.csv.gzip",
	"https://placedata.reddit.com/data/canvas-history/2022_place_canvas_history-000000000022.csv.gzip",
	"https://placedata.reddit.com/data/canvas-history/2022_place_canvas_history-000000000023.csv.gzip",
	"https://placedata.reddit.com/data/canvas-history/2022_place_canvas_history-000000000024.csv.gzip",
	"https://placedata.reddit.com/data/canvas-history/2022_place_canvas_history-000000000025.csv.gzip",
	"https://placedata.reddit.com/data/canvas-history/2022_place_canvas_history-000000000026.csv.gzip",
	"https://placedata.reddit.com/data/canvas-history/2022_place_canvas_history-000000000027.csv.gzip",
	"https://placedata.reddit.com/data/canvas-history/2022_place_canvas_history-000000000028.csv.gzip",
	"https://placedata.reddit.com/data/canvas-history/2022_place_canvas_history-000000000029.csv.gzip",
	"https://placedata.reddit.com/data/canvas-history/2022_place_canvas_history-000000000030.csv.gzip",
	"https://placedata.reddit.com/data/canvas-history/2022_place_canvas_history-000000000031.csv.gzip",
	"https://placedata.reddit.com/data/canvas-history/2022_place_canvas_history-000000000032.csv.gzip",
	"https://placedata.reddit.com/data/canvas-history/2022_place_canvas_history-000000000033.csv.gzip",
	"https://placedata.reddit.com/data/canvas-history/2022_place_canvas_history-000000000034.csv.gzip",
	"https://placedata.reddit.com/data/canvas-history/2022_place_canvas_history-000000000035.csv.gzip",
	"https://placedata.reddit.com/data/canvas-history/2022_place_canvas_history-000000000036.csv.gzip",
	"https://placedata.reddit.com/data/canvas-history/2022_place_canvas_history-000000000037.csv.gzip",
	"https://placedata.reddit.com/data/canvas-history/2022_place_canvas_history-000000000038.csv.gzip",
	"https://placedata.reddit.com/data/canvas-history/2022_place_canvas_history-000000000039.csv.gzip",
	"https://placedata.reddit.com/data/canvas-history/2022_place_canvas_history-000000000040.csv.gzip",
	"https://placedata.reddit.com/data/canvas-history/2022_place_canvas_history-000000000041.csv.gzip",
	"https://placedata.reddit.com/data/canvas-history/2022_place_canvas_history-000000000042.csv.gzip",
	"https://placedata.reddit.com/data/canvas-history/2022_place_canvas_history-000000000043.csv.gzip",
	"https://placedata.reddit.com/data/canvas-history/2022_place_canvas_history-000000000044.csv.gzip",
	"https://placedata.reddit.com/data/canvas-history/2022_place_canvas_history-000000000045.csv.gzip",
	"https://placedata.reddit.com/data/canvas-history/2022_place_canvas_history-000000000046.csv.gzip",
	"https://placedata.reddit.com/data/canvas-history/2022_place_canvas_history-000000000047.csv.gzip",
	"https://placedata.reddit.com/data/canvas-history/2022_place_canvas_history-000000000048.csv.gzip",
	"https://placedata.reddit.com/data/canvas-history/2022_place_canvas_history-000000000049.csv.gzip",
	"https://placedata.reddit.com/data/canvas-history/2022_place_canvas_history-000000000050.csv.gzip",
	"https://placedata.reddit.com/data/canvas-history/2022_place_canvas_history-000000000051.csv.gzip",
	"https://placedata.reddit.com/data/canvas-history/2022_place_canvas_history-000000000052.csv.gzip",
	"https://placedata.reddit.com/data/canvas-history/2022_place_canvas_history-000000000053.csv.gzip",
	"https://placedata.reddit.com/data/canvas-history/2022_place_canvas_history-000000000054.csv.gzip",
	"https://placedata.reddit.com/data/canvas-history/2022_place_canvas_history-000000000055.csv.gzip",
	"https://placedata.reddit.com/data/canvas-history/2022_place_canvas_history-000000000056.csv.gzip",
	"https://placedata.reddit.com/data/canvas-history/2022_place_canvas_history-000000000057.csv.gzip",
	"https://placedata.reddit.com/data/canvas-history/2022_place_canvas_history-000000000058.csv.gzip",
	"https://placedata.reddit.com/data/canvas-history/2022_place_canvas_history-000000000059.csv.gzip",
	"https://placedata.reddit.com/data/canvas-history/2022_place_canvas_history-000000000060.csv.gzip",
	"https://placedata.reddit.com/data/canvas-history/2022_place_canvas_history-000000000061.csv.gzip",
	"https://placedata.reddit.com/data/canvas-history/2022_place_canvas_history-000000000062.csv.gzip",
	"https://placedata.reddit.com/data/canvas-history/2022_place_canvas_history-000000000063.csv.gzip",
	"https://placedata.reddit.com/data/canvas-history/2022_place_canvas_history-000000000064.csv.gzip",
	"https://placedata.reddit.com/data/canvas-history/2022_place_canvas_history-000000000065.csv.gzip",
	"https://placedata.reddit.com/data/canvas-history/2022_place_canvas_history-000000000066.csv.gzip",
	"https://placedata.reddit.com/data/canvas-history/2022_place_canvas_history-000000000067.csv.gzip",
	"https://placedata.reddit.com/data/canvas-history/2022_place_canvas_history-000000000068.csv.gzip",
	"https://placedata.reddit.com/data/canvas-history/2022_place_canvas_history-000000000069.csv.gzip",
	"https://placedata.reddit.com/data/canvas-history/2022_place_canvas_history-000000000070.csv.gzip",
	"https://placedata.reddit.com/data/canvas-history/2022_place_canvas_history-000000000071.csv.gzip",
	"https://placedata.reddit.com/data/canvas-history/2022_place_canvas_history-000000000072.csv.gzip",
	"https://placedata.reddit.com/data/canvas-history/2022_place_canvas_history-000000000073.csv.gzip",
	"https://placedata.reddit.com/data/canvas-history/2022_place_canvas_history-000000000074.csv.gzip",
	"https://placedata.reddit.com/data/canvas-history/2022_place_canvas_history-000000000075.csv.gzip",
	"https://placedata.reddit.com/data/canvas-history/2022_place_canvas_history-000000000076.csv.gzip",
	"https://placedata.reddit.com/data/canvas-history/2022_place_canvas_history-000000000077.csv.gzip",
	"https://placedata.reddit.com/data/canvas-history/2022_place_canvas_history-000000000078.csv.gzip",
}

func DownloadAll() {
	fmt.Println("Download Tool initialized")
	l := len(csvParts)
	dl := make(chan string, l)

	for _, link := range csvParts {
		dl <- link
	}
	close(dl)

	fmt.Printf("(%d) parts to be downloaded, please be patient and do not stop the process\n", l)
	var wg sync.WaitGroup
	for i := 0; i < MaxConcurrentDownloads; i++ {
		wg.Add(1)
		go func() {
			for link := range dl {
				fmt.Printf("-> downloading %s\n...", link)
				resp, err := http.Get(link)
				if err != nil {
					panic(err)
				}

				func() {
					defer resp.Body.Close()

					status := resp.StatusCode
					if status != 200 {
						panic(
							fmt.Sprintf(
								"download error: expected 200 status but got %d",
								status,
							))
					}

					f, err := os.Create(filenameFromURL(link))
					if err != nil {
						panic(err)
					}
					defer f.Close()

					_, err = io.Copy(f, resp.Body)
					if err != nil {
						panic(err)
					}
					fmt.Printf("<- downloaded %s\n", link)
				}()
			}
			wg.Done()
		}()
	}
	wg.Wait()
}

func filenameFromURL(rawurl string) string {
	url, err := url.Parse(rawurl)
	if err != nil {
		panic(err)
	}
	parts := strings.Split(url.Path, "/")
	if len(parts) == 0 {
		panic("unexpected URL length, check csvPart links")
	}
	return parts[len(parts)-1]
}
