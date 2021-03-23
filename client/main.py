#!/usr/bin/env python

import client
import websockets
import asyncio

async def register():
    uri = "ws://localhost:8001"
    async with websockets.connect(uri) as websocket:
        
        c = client.Client("user1")
        c.show_keys()
        name = c.username

        await websocket.send(c.create_register_msg())

        response = await websocket.recv()
        print(f"< {response}")


asyncio.get_event_loop().run_until_complete(register())

