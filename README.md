# suspx - <img src="https://user-images.githubusercontent.com/22831717/162569206-6e865d46-9c42-4a1e-9bbe-a72c9415e4ca.png" width="50">
Analytical tool for sus pixels in r/place.

## Introduction

I built this tool to analyze datasets from [r/Place](https://reddit.com/r/place/). The main goal was to detect and filter genuine human interactions 
out of the canvas to obtain a canvas representation of non-human interactions. But because it works by running simulations, pixel by pixel in order, 
it could be re-factorized and adapted to other use cases to achieve a single tool for multiple canvas analysis and save time by avoiding having to code
python scripts for every analysis of future r/Place events.

The tool also includes a download tool so you don't have to download the CSVs separately (but you still have to decompress them by your own before
running the simulations). It is also important to point out that the CSV datasets are not in order (0 may not be the first part and 78 may not be the
last one), so this tool will help you by reading the first row of each dataset and ordering them automatically. All you need is to have all the CSV 
datasets downloaded in the same directory before running the simulation.

## How does it work?

It runs a simulation with provided datasets (CSV files), reading and processing pixel by pixel in the corresponding order. A hashtable stores a pixel
bucket for each user (hashed) name. Each bucket contains the needed information for later analysis: sus pixels count (suspicious pixels), pixels placed,
last pixel placed.

Regarding bot analysis, a sus pixel is added to the bucket depending on the **time margin** (m) and **cooldown** (cd) parameters. First a **cooldown** is set 
(5 min), the **time margin** determines the extra time added to the **cooldown** (basically, the margin gives extra time to the user to react since the 
next pixel is available). If the delta (Δ) of the current pixel timestamp being processed - last pixel by the same user timestamp is less than this time, 
the pixel is considered suspicious and added to the bucket. 

The **threshold** (t) parameter determines how many consecutive pixels (a break of time equal or greater than cooldown+margin breaks the streak and resets 
the suspicious counter in the bucket) by the same user are needed for the pixels to be drawn on the canvas. If the sus count reaches this **threshold**, 
all the following suspicious pixels will be drawn on the canvas until a non-suspicious pixel is found.

## Finding optimal parameters

At a glance, a short **time margin** seems like a good idea, a low time of reaction is very non-human. And this is true, some pixels will be rendered and
you can assume with a fair ammount certainty that they are non-human, but a low number of pixels will be displayed. This is because most bots actually
delay their next pixels or even randomize them to ensure availability and mitigate detection.

I took the most popular bot as an example: https://github.com/Skeeww/Bot and https://github.com/PlaceNL/Bot, which are basically the same forked version
of a public bot that many communities use. It works by using a extension in the client which connects (via websockets) to a command server that coordinate
and sends orders. In the following line: 
https://github.com/Skeeww/Bot/blob/f131b29960544e1f8123b89f50aa7c903767dc4f/bot.js#L308 we can see how the `next pixel = <cooldown> + 3s + 
<random number between 0 and 10>`. That's why I set the default time margin to 14010 ms. A time span of 14s will give many false positives for low 
thresholds, but as the threshold grows, it yields a quite representative canvas of non-human interactions (since very few humans will place pixels for 1h,
2h, 4h, 8h, 10h or 12h in a consistent way, maybe some will, but that is also represented by the intensity of the colors that depends on the number of
pixels placed).

I've also tried simulations with greater **cooldown**'s, 20 minutes for example for unverified accounts. But it will overlap with lower **cooldowns** 
since a 20 min cooldown with 1 threshold, will also include 5 min cooldowns with 4 sus pixels. So on the lower side of **threshold** it won't yield 
representative results and with larger **threshold**s the results are similar to the ones with cd = 5.

So that leaves us with a `m = 14010` ms, `cd = 5` minutes and then run it with different **thresholds**: `6 (=30min), 12 (=1h), 24 (=2h), 48 (4=h), 
96 (=8h), 100 (=10h), 144 (=12h), etc.` for a cd of 5.

But please, try running it with different parameters with different reasoning and show us!

### Hints

Use document and proven cases of bots as a reference that you are on the right track while adjusting parameters. 

- The BTS logo in the bottom left corner is botted by spanish streamers during a short time of period (below 2h) on a live stream on twitch due to a 
misunderstanding in a french vs spanish 'war': https://clips.twitch.tv/RepleteExuberantWolfAsianGlow-3KNYwHYxJxGLFDJp
- r/PlaceNL used the same bot with a great number of clients (around 2000-3000) and the command server orders were public here: https://placenl.noahvdaa.me/. 
A screenshot of the site: ![Screen Shot 2022-04-09 at 14 57 28](https://user-images.githubusercontent.com/22831717/162575514-77b82728-642a-4186-9ff7-793451efb3ad.png)

## Usage

1. Download the tool executable (in releases) or compile it by your own. 

Minimum 16GB RAM is recommended since the hashtable will grow to roughly 8GB 
containing all the users who placed at least 1 pixel. You will also need about 21GB or more of free space for the decompressed datasets. A SSD is 
recommended to speed up the process and a good internet connection to download all the parts.

2. Execute the tool with the -d parameter to download all the datasets: `./suspx -d`. Please be patient and check that you have all the parts. In the
event of a missing part or a problem, delete the datasets and try it again.

3. Decompress the parts. In linux you may want to use: `gunzip *.gzip` or `gzip -d *.gzip`. If you get an error, try renaming them from .gzip to .gz, 
here is a one liner command: `for f in *.gzip; do mv -- "$f" "${f%.gzip}.gz"; done`. 

On Windows or MacOS you just have to find an appropiate zip tool for gzip files. 

Again, please verify that all the CSV parts are available in the same directory as the executable. From 0 to 78.

4. Run the tool with the desired parameters. See available parameters below or type `./suspx --help` in your terminal.

## Parameters

- `-cd <minutes>`. Time in minutes for the cooldown. (default: 5)
- `-threshold <number>`. Suspicious threshold, above this threshold of consecutive pixels (or non-consecutives if -nc is passed down), the following 
consecutive pixels (or not) will be drawn (default: 12)
- `-margin <milliseconds> or `-m <milliseconds>`. Time margin that will be added to the cooldown (default: 14010).
- `-d`. Run the download tool instead and exits.
- `-h <size>`. Pixel height of the canvas. It defaults to the corresponding size of 2022 r/Place.
- `-w <size>`. Pixel width of the canvas. It defaults to the corresponding size of 2022 r/Place.
- `-nc`. Run the tool in non-consecutive mode. All suspicious will be drawn on the canvas, no streaks required. Only the current pixel and the last 
  one is considered.
- `-o <name.png>`. Provide a different name for the resulting exported PNG file (default: 'res.png')
  
## Results
  
  ### Threshold=6
  <img src="https://github.com/pmrt/suspx/raw/master/results/cd5_m14010_t6.png" width="500">
  
   ### Threshold=12
   <img src="https://github.com/pmrt/suspx/raw/master/results/cd5_m14010_t12.png" width="500">
  
   ### Threshold=24
   <img src="https://github.com/pmrt/suspx/raw/master/results/cd5_m14010_t24.png" width="500">
  
   ### Threshold=48
   <img src="https://github.com/pmrt/suspx/raw/master/results/cd5_m14010_t48.png" width="500">
  
   ### Threshold=96
   <img src="https://github.com/pmrt/suspx/raw/master/results/cd5_m14010_t96.png" width="500">
  
   ### Threshold=100
   <img src="https://github.com/pmrt/suspx/raw/master/results/cd5_m14010_t100.png" width="500">
  
   ### Threshold=144
   <img src="https://github.com/pmrt/suspx/raw/master/results/cd5_m14010_t144.png" width="500">


## Contributors
  
- Pedro M. M. <s+gh@pedro⚫️to>; where ⚫️ = '.'

That's all, find the sus pixels!
