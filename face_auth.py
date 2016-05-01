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

detectURL = "https://api.projectoxford.ai/face/v1.0/detect?returnFaceId=true&returnFaceLandmarks=false&returnFaceAttributes=age"
detectHeaders = {
	"Content-Type": "application/json",
	"Ocp-Apim-Subscription-Key": "38c44ac804c44f6e97673d815163a1db"
}

res = requests.post(detectURL, data=json.dumps({"url":"http://hellogiggles.com/wp-content/uploads/2015/04/10/anna-kendrick-pitch-perfect-650-430.jpg"}), headers=detectHeaders)
# prints response
#print(res.content)

# prints faceId
detectDict = json.loads(res.content)[0]
print(detectDict["faceId"])








########################
NOT USEFUL BELOW HERE
########################



# file object.p
#DATA = open("facescrub_actors.txt", "r")

#nameSet = set()
# def checkValid(arg):
# 	if arg[0] in nameSet:
# 		# we already had URL for this name
# 		return False
# 	try:
# 		r = requests.get(arg[3])
# 		#r = requests.get("asdfasdfdsfsafsaf.com")
# 		nameSet.add(arg[0])
# 		return r.status_code == 200
# 	except requests.exceptions.RequestException:
# 		return False

# # split on tabs
# res = map((lambda x: x.split("\t")), DATA)

# # filter out invalid URL
# valid = filter(checkValid, res)

# #for data in valid:
# #	print(data[3])


# count = 0
# retArray = []

# #nameSet = set()
# for data in valid:
# 	#if data[0] not in nameSet:
# 		newDictionary = {}
# 		newDictionary["name"] = data[0]
# 		newDictionary["url"] = data[3]
# 		newDictionary["id"] = count

# 		retArray.append(newDictionary)
# 		#nameSet.add(data[0])
# 		count += 1

# #print(repr(retArray))

# retJSON = json.dumps(retArray, sort_keys=True, indent=4, separators=(",", ": "))
# #print(retJSON)
# f = open("faces.json", "w")
# f.write(retJSON)









# # c = httplib.HTTPConnection('www.example.com')
# # c.request("HEAD", '')
# # if c.getresponse().status == 200:
# #    print('web site exists')


# # a = [1,2,3,4]

# # def add1(arg):
# # 	return arg+1

# # out = map(add1, a)

# # for b in out:
# # 	print(str(b))