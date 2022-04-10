# suspx - <img src="https://user-images.githubusercontent.com/22831717/162569206-6e865d46-9c42-4a1e-9bbe-a72c9415e4ca.png" width="50">
Analytical tool for sus pixels in r/place.

## Table of Contents

- [1. Introduction](#1-introduction)
- [2. Features](#2-features)
- [3. How does it work](#3-how-does-it-work)
- [4. Instrument example, how does the bot analysis work?](#4-instrument-example-how-does-the-bot-analysis-work)
  - [4.1. Finding optimal parameters for the bot analysis](#41-finding-optimal-parameters-for-the-bot-analysis)
  - [4.2. Hints](#42-hints)
- [5. Usage](#5-usage)
- [6. Parameters](#6-parameters)
- [7. Results](#7-results)
  - [7.1. Bot instrument results](#71-bot-instrument-results)
- [8. Instrument API](#8-instrument-api)
- [9. Contributors](#9-contributors)


## 1. Introduction

I built this tool to analyze datasets from
[r/Place](https://reddit.com/r/place/). The main goal was to detect and filter
genuine human interactions out of the canvas to obtain a canvas representation
of non-human interactions. But because it works by running simulations, pixel by
pixel in the correct order, it is very easy to add other use cases (see the
instruments API below) to achieve a single tool for multiple r/Place analysis,
saving a considerable amount of time.

## 2. Features

- Automatically download all the needed CSV files
- Automatically order all the parts in the correct chronological order
- An instrument API to easily create and adapt it to new analysis and use cases

## 3. How does it work?

It runs a simulation with the datasets (CSV files) found in the same path as the
executable, reading and processing pixel by pixel in the correspondig order. The
different analysis tools are called **instruments**. Each placed pixel involves
an **instrument bucket** being added/retrieved to/from a hashtable grouped by
users, so the **bucket** is shared across different pixels iterations for the
same user. Each **instrument** decides what to store in its bucket, for example
the bot instrument stores that last pixel placed by the same user and the number
of suspicious pixels.

On each pixel iteration, the **instrument** is invoked, passing down the pixel
details, the bucket for the corresponding user and the whole hashtable if
needed. When the simulation process ends, the **instrument** is called again
passing down the resulting hashtable and canvas.

## 4. Instrument example, how does the bot analysis work?

Regarding bot analysis, the bot **instrument** defines a set of arguments to be
passed to the program: a **time margin** (m) in milliseconds, the **cooldown**
(cd) in minutes and a **threshold** (t).

- The **cooldown** determines the cooldown between pixels (default: 5 minutes).
- The **time margin** determines the extra time added to the **cooldown**, in
  order words, the margin defines the extra time for the user to react since the
  next pixel is available. Once this margin is surpassed, the pixel is no longer
  considered suspicious.
- The **threshold** parameter determines how many consecutive suspicious pixels
  by the same user are needed for the following suspicious pixels to be drawn on
  the canvas. When the suspicious count in the bucket of a given user reaches
  this threshold, all the following suspicious pixels will be drawn on the
  canvas until a non-suspicious pixel is found, when the suspicious count in the
  bucket for the corresponding user will be reset.

Considering `Δ = Tpx1 - Tpx`, let `Tpx1` be the timestamp of the current pixel
being processed and let `Tpx` be the timestamp of the last pixel by the same
user, the current pixel is considered suspicious if `Δ < cd + cd`.

So, following these principles, the bot instruments defines a **instrument
bucket** where it stores the suspicious count and the last pixel for each new
user. Then on each pixel iteration the bot instrument is invoked and it checks
the condition above in the corresponding bucket. If the condition is met and the
threshold is reached, the bot instruments instructs the simulator to draw the
following suspicious pixels on the canvas until a non-suspicious pixel is found.

### 4.1. Finding optimal parameters for the bot analysis
At first, the naive approach seems pretty straightforward: a short **time
margin** seems like a good idea because a low time of reaction is very
non-human. And that is true, some pixels will be rendered and you can assume
with a fair ammount certainty that they are non-human, but you will find
yourself with a very low number of pixels being drawn on the canvas. This is
because most bots actually delay their next pixels or even randomize them to
ensure availability and mitigate detection.

I took the most popular bot as an example: https://github.com/Skeeww/Bot and
https://github.com/PlaceNL/Bot, which are basically the same forked version of a
public bot that many communities use. It works by using a extension in the
client which connects (via websockets) to a command server that coordinates and
sends orders. In the following line:
https://github.com/Skeeww/Bot/blob/f131b29960544e1f8123b89f50aa7c903767dc4f/bot.js#L308
we can see how the `next pixel = <cooldown> + 3s + <random number between 0 and
10>`. That's why I set the default time margin to 14010 ms. A time span of 14s
will give many false positives for low thresholds, but as the threshold grows,
it yields quite a good canvas representation of non-human interactions (since
very few humans will place pixels for 1h, 2h, 4h, 8h, 10h or 12h in a consistent
way, maybe some will, but that is also depicted by the intensity of the colors
that depends on the number of pixels placed).

I've also tried simulations with greater **cooldown**s, 20 minutes for example
for unverified accounts. But it will overlap with lower **cooldowns** since a 20
  min cooldown with 1 threshold, will also include 5 min cooldowns with 4 sus
  pixels. So on the lower side of **threshold** it won't yield representative
  results and with larger **threshold**s the results are similar to the ones
  with cd = 5.

That leaves us with a `m = 14010` ms, `cd = 5` minutes and we just need to run
it with different **thresholds**: `6 (=30min), 12 (=1h), 24 (=2h), 48 (4=h), 96
(=8h), 100 (=10h), 144 (=12h), etc. (for a cd of 5min)`.

But please, try running it with different parameters with different reasoning
and show us!

### 4.2. Hints

Use document and proven cases of bots as a reference that you are on the right
track while adjusting parameters.

- The BTS logo in the bottom left corner is botted by spanish streamers during a
  short time of period (below 2h) on a live stream on twitch due to a
  misunderstanding in a french vs spanish 'war':
  https://clips.twitch.tv/RepleteExuberantWolfAsianGlow-3KNYwHYxJxGLFDJp
- r/PlaceNL used the same bot with a great number of clients (around 2000-3000)
  and the command server orders were public here: https://placenl.noahvdaa.me/.
A screenshot of the site: ![Screen Shot 2022-04-09 at 14 57 28](https://user-images.githubusercontent.com/22831717/162575514-77b82728-642a-4186-9ff7-793451efb3ad.png)

## 5. Usage

1. Download the tool executable (in releases) or compile it by your own.

Minimum 16GB RAM is recommended since the hashtable will grow to roughly 8GB for
the 2022 version, containing all the users who placed at least 1 pixel. You will
also need about 21GB or more of free space for the decompressed datasets. A SSD
is recommended to speed up the process and a good internet connection to
download all the parts.

2. Execute the tool with no datasets. It will ask you if you want it to download
   all the parts for you, type `y` or `yes` and wait. Please, be patient and
   check that you have all the parts. In the case of a missing part or a
   problem, delete the datasets involved and try it again.

3. Decompress the parts. In linux you may want to use: `gunzip *.gzip` or `gzip
   -d *.gzip`. If you get an error, try renaming them from .gzip to .gz, here is
   a one liner command: `for f in *.gzip; do mv -- "$f" "${f%.gzip}.gz"; done`.
   On Windows or MacOS you just have to find an appropiate zip tool for gzip
   files.

Again, please verify that all the CSV parts are available in the same directory
as the executable. From 0 to 78.

4. Run the tool with the desired parameters. See available parameters below or
   type `./suspx --help` in your terminal.

5. Once you have all the parts extracted in the same directory and everything is
   ready, it will ask you to select an instrument to be executed during the
   simulation. Pick one from the list.

## 6. Parameters

### General
- `-h <size>`. Pixel height of the canvas. It defaults to the corresponding size of 2022 r/Place.
- `-w <size>`. Pixel width of the canvas. It defaults to the corresponding size of 2022 r/Place.
- `-o <name.png>`. Provide a different name for the resulting exported PNG file (default: 'res.png')

### Bot instrument

- `-cd <minutes>`. Time in minutes for the cooldown. (default: 5)
- `-threshold <number>`. Suspicious threshold, above this threshold of
  consecutive pixels (or non-consecutives if -nc is passed down), the following
  consecutive pixels (or not) will be drawn (default: 12)
- `-margin <milliseconds> or `-m <milliseconds>`. Time margin that will be added
  to the cooldown (default: 14010).
- `-nc`. Run the tool in non-consecutive mode. All suspicious will be drawn on
  the canvas, no streaks required. Only the current pixel and the last one is
  considered.

## 7. Results

### 7.1. Bot instrument results:

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

## 8. Instrument API

TODO - Describe the exact methods for the instruments and buckets. In the
meantime you can take a look at the instruments in the `instruments` package if
you want to add your own instrument for a new analysis.

## 9. Contributors

- Pedro M. M. <s+gh@pedro⚫️to>; where ⚫️ = '.'

That's all, find the sus pixels!
