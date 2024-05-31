import time

import pika
import json

connection = pika.BlockingConnection()
channel = connection.channel()

node_name = 'repeater'

channel.queue_declare(queue=node_name)
channel.queue_declare(queue=node_name+'_response')

print("Ready to handle messages")
while True:
    method_frame, properties, body = next(channel.consume(node_name))
    channel.basic_ack(method_frame.delivery_tag)

    data = json.loads(body.decode("utf-8"))
    print("Request: ", data)

    answer = bytes(json.dumps(data), "utf-8")
    print("Sending data: ", answer)
    channel.basic_publish(exchange='', routing_key=node_name + '_response', body=answer)
