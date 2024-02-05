import json 
from random import  randrange 

from locust import HttpUser, between, task  

debug = False

def printDebug(msg):
    if debug: 
        print(msg)


class Reader():

    def __init__(self):
        self.array = []

    def pickRandom(self):    
        lenght = len(self.array)

        if(lenght > 0):
            random_index = randrange(0, lenght - 1) if lenght > 1 else 0 
            return self.array.pop(random_index)
        
        else:
            print(">> Reader: We can't found a register")
            return None

    def load(self):
        print(">> Reader: Loading.....")

        try:
         with open("data-traffic.json", 'r')  as data_file:
             self.array = json.loads(data_file.read())

        except Exception as error:
            print(f'Error loading data from file: {error}' )

class MessaggeTraffic(HttpUser):
    wait_time = between(0.1, 0.9) # Wait time between requests in seconds (defaults to between(1,2))
    reader = Reader()
    reader.load()

    def on_start(self):
        print(">>MessaggeTraffic: Sending traffic...")
    
    @task
    def PostMessagge(self):
        random_data = self.reader.pickRandom()
        if(random_data != None):
            data_to_send = json.dumps(random_data)
            printDebug(data_to_send)

            self.client.post("/", json=random_data)
        else:
            print(">>MessaggeTraffic: Traffic sender end...")
            self.stop(True)

    @task
    def GetMessage(self):
        self.client.get("/")
        
