#!/usr/bin/python3

# double nested array
# n number of fixed size things
# json 
# http request to get the file
# curl everything, remove bad links
# put into JSON
#

#print("hello world")

#import urllib.request
#import http.client
import requests
import json
#r = urllib.request.urlopen("https://www.marktai.com/upload/facescrub_actors.txt").read()


# file object.p
DATA = open("facescrub_actors_test.txt", "r")

#for line in DATA:
	#print(line)
#	ret = line.split("\t")
#	print(ret[3])
#print(DATA.readline())

def checkValid(arg):
	#print(arg[3])
	#return True
	#r = http.client.HTTPConnection("http://upload.wikimedia.org/wikipedia/commons/5/5d/AaronEckhart10TIFF.jpg", 80)
	#r.request("HEAD", '')
	#r.connect()
	#r.request("GET", "/")
	#return r.getresponse().status == 200
	#return True
	#print(arg[3][1])
	#return True
	try:
		r = requests.get(arg[3])
		#r = requests.get("asdfasdfdsfsafsaf.com")
		return r.status_code == 200
	except requests.exceptions.RequestException:
		return False

# split on tabs
res = map((lambda x: x.split("\t")), DATA)

# filter out invalid URL
valid = filter(checkValid, res)

#for data in valid:
#	print(data[3])


count = 0
retArray = []

nameSet = set()
for data in valid:
	if data[0] not in nameSet:
		newDictionary = {}
		newDictionary["name"] = data[0]
		newDictionary["url"] = data[3]
		newDictionary["id"] = count

		retArray.append(newDictionary)
		nameSet.add(data[0])
		count += 1

#print(repr(retArray))

retJSON = json.dumps(retArray)
f = open("faces.json", "w")
f.write(retJSON)









# c = httplib.HTTPConnection('www.example.com')
# c.request("HEAD", '')
# if c.getresponse().status == 200:
#    print('web site exists')


# a = [1,2,3,4]

# def add1(arg):
# 	return arg+1

# out = map(add1, a)

# for b in out:
# 	print(str(b))