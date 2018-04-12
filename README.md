# AQI-Monitoring

Get AQI (Air Quality Index) data from [行政院環境保護署](https://www.epa.gov.tw/) and push notify to user when air quality is 'GOOD'.

# How to use it ?

First of all , Add this [LINE robot](https://line.me/R/ti/p/%40zhm1865k) as friend and activate LINE Notify by your own custom activate url.

![螢幕快照 2018-04-12 下午4.04.38](https://i.imgur.com/c18t4TN.png)

After that , you can register observatory nearby .

Send your location to server and we will find out which observatory is nearest.

![螢幕快照 2018-04-12 下午4.04.38](https://i.imgur.com/q7Jq9WK.png)

Of course , you can register observatory as much as you want (your home, workplace , …).

![螢幕快照 2018-04-12 下午4.04.38](https://i.imgur.com/v0SOXlw.png)

Now , when your observatory shows AQI is under 51 (which means air quality is good in this area) , you will get notify and a adorable Pokémon pic randomly.

(P.S. Notify only push  between 8 am and 9 am , 6 pm and 7 pm)

![螢幕快照 2018-04-12 下午4.04.38](https://i.imgur.com/Whm9bZG.png)

# Others

#### Calculate distance of two location

```
func distance(lat1, lon1, lat2, lon2 float64) float64 {
	// convert to radians
	// must cast radius as float to multiply later
	var la1, lo1, la2, lo2, r float64
	la1 = lat1 * math.Pi / 180
	lo1 = lon1 * math.Pi / 180
	la2 = lat2 * math.Pi / 180
	lo2 = lon2 * math.Pi / 180

	r = 6378100 // Earth radius in METERS

	// calculate
	h := hsin(la2-la1) + math.Cos(la1)*math.Cos(la2)*hsin(lo2-lo1)

	return 2 * r * math.Asin(math.Sqrt(h))
}
```

```
func hsin(theta float64) float64 {
	return math.Pow(math.Sin(theta/2), 2)
}
```

#### Line Message API plan

![螢幕快照 2018-04-12 下午4.04.38](https://i.imgur.com/ay8pUgZ.png)

Although the free plan cant't push message , this plan can get unlimited number of friends.

#### Line Notify API

Then, use Line Notify API to implement AQI notify function.

![螢幕快照 2018-04-12 下午4.04.38](https://i.imgur.com/pucTx8a.png)

