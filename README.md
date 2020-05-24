# mtbsearch

A small project to index and search through Mountain Bikes on https://www.vitalmtb.com/
based on specifications of the bikes.

## Usage

1. Run the fetch script to download the specs: `go run cmd/fetch/main.go`
1. Run the index script to index the specs: `go run cmd/index/main.go`
1. Run the search script to search the index: `go run cmd/search/main.go`

Modify the search script to change the search parameters and output.

## Example

At the time of writing, the following Bleve search will yield the results below:

```go
stringQuery := strings.Join([]string{
	"+wheel_size_inches:29",
	"+model_year_float:>=2019",
	"+brakes:shimano",
	"+shifters:shimano",
	"+rear_travel_mm:>=130",
	"+rear_travel_mm:<=150",
	"+fork_travel_mm:>=150",
	"+fork_travel_mm:<=160",
	"+bottle_cage_mounts:yes",
	"+frame_material:carbon",
	"-e_bike_class:class",
}, " ")
```

```
$ go run cmd/search/main.go 
2020/05/24 17:23:36 2020 Ibis Ripmo V2 Carbon XTR Bike - $9,299
2020/05/24 17:23:36 2020 Ibis Ripmo V2 Carbon XT Bike - $5,899
2020/05/24 17:23:36 2020 Ibis Ripmo V2 Carbon SLX Bike - $5,199
2020/05/24 17:23:36 2019 Eminent Onset LT Advanced Bike - $5,199
2020/05/24 17:23:36 2020 Orbea Occam M-LTD Bike - $7,999
2020/05/24 17:23:36 2020 Yeti SB130 T1 Bike - Standard:$6,799With XMC wheelset:$8,299
2020/05/24 17:23:36 2020 Evil Offering XT Bike - $5,799
2020/05/24 17:23:36 2020 Evil Offering XTR Bike - $7,399
2020/05/24 17:23:36 2020 Ibis Ripmo Carbon XTR Bike - $9,299
2020/05/24 17:23:36 2020 Ibis Ripmo Carbon XT Bike - $5,899
2020/05/24 17:23:36 2020 Ibis Ripmo Carbon SLX Bike - $5,199
2020/05/24 17:23:36 2020 Alchemy Arktos 29 XTR Bike - $9,499
2020/05/24 17:23:36 2020 Santa Cruz Hightower Carbon CC XTR Reserve Bike - $9,899
2020/05/24 17:23:36 2020 Trek Slash 9.9 XTR Bike - $7,499.99
2020/05/24 17:23:36 2020 Trek Slash 9.8 XT Bike - $5,899.99
2020/05/24 17:23:36 2019 Eminent Onset LT Pro Bike - $7,599
2020/05/24 17:23:36 2019 Evil Offering XT Jenson USA Exclusive Bike - $5,599.99
2020/05/24 17:23:36 2020 Alchemy Arktos 29 XT 12 Speed Bike - $6,099
2020/05/24 17:23:36 2020 Mondraker Foxy Carbon RR Bike - N/A
2020/05/24 17:23:36 2020 Scott Genius 910 Bike - $4,999.99
2020/05/24 17:23:36 2020 KTM Prowler Sonic Bike - N/A
2020/05/24 17:23:36 2020 Cube Stereo 150 C:62 SL 29 Bike - N/A
2020/05/24 17:23:36 2020 KTM Prowler Glory Bike - N/A
2020/05/24 17:23:36 2020 Intense Primer 29" XTR Jenson USA Exclusive Bike - $8,400
2020/05/24 17:23:36 2020 Niner RIP 9 RDO 29 5-Star Shimano XTR LTD Bike - $9,100
2020/05/24 17:23:36 2020 Specialized Stumpjumper S-Works 29 Bike - $9,520
2020/05/24 17:23:36 2020 Canyon Strive CF 8.0 Bike - $4,699
2020/05/24 17:23:36 2019 Santa Cruz Hightower LT Carbon CC XTR Reserve Bike - $9,599
2020/05/24 17:23:36 2020 Niner RIP 9 RDO 29 4-Star Shimano XT Bike - $6,600
2020/05/24 17:23:36 2020 Niner RIP 9 RDO 29 3-Star Shimano XT Bike - $5,700
2020/05/24 17:23:36 2019 Cube Stereo 150 C:62 Race 29 Bike - N/A
2020/05/24 17:23:36 2019 Cube Stereo 150 C:68 Action Team 29 Bike - N/A
2020/05/24 17:23:36 2020 Pivot Switchblade Pro XT/XTR 29 - With DT Swiss wheels:$6,799With Reynolds wheels upgrade option:$8,099
2020/05/24 17:23:36 2020 Pivot Switchblade Pro XT/XTR FOX Live Valve 29 Bike - With DT Swiss wheels:$8,699With Reynolds wheels upgrade option:$9,999
2020/05/24 17:23:36 2020 Pivot Switchblade Team XTR 29 Bike - $8,999
2020/05/24 17:23:36 2020 Pivot Switchblade Team XTR FOX Live Valve 29 Bike - $10,899
2020/05/24 17:23:36 2020 Pivot Switchblade Race XT 29 Bike - $5,499
```

## Notes

 - Specs data is stored at `data.json`
 - Index is stored at `index.bleve`
 - Bleve query string docs can be found at https://blevesearch.com/docs/Query-String-Query/