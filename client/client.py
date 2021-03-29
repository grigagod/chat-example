import json
from coincurve import PrivateKey
import binascii


class Client:
    def __init__(self, username):
        self.username = username
        self.init_keys()

    def init_keys(self):
        f = open(self.username + '.json')

        data = json.load(f)

        priv_key_hex = data["private_key"]

        self.priv_key = PrivateKey.from_hex(priv_key_hex)

        f.close()

    def show_keys(self):
        print(self.priv_key.to_hex())
        print(binascii.hexlify(self.priv_key.public_key.format()))


    def create_register_msg(self):
        return json.dumps({'Type' : 2, 
                          'username' : self.username, 
                          'pub_key' : binascii.hexlify(self.priv_key.public_key.format()).decode("utf-8")})        
