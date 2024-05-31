import time

import pika
import json

connection = pika.BlockingConnection()
channel = connection.channel()

node_name = 'task_splitter_formatter'

channel.queue_declare(queue=node_name)
channel.queue_declare(queue=node_name+'_response')

print("Ready to handle messages")
while True:
    method_frame, properties, body = next(channel.consume(node_name))
    channel.basic_ack(method_frame.delivery_tag)

    data = json.loads(body.decode("utf-8"))
    print("Request: ", data)

    system_prompt = ""
    prompt = ""
    for input_wrapper in data:
        if input_wrapper['wrapper_id'] == "system_prompt":
            system_prompt = input_wrapper["value_info"]
        if input_wrapper['wrapper_id'] == "prompt":
            prompt = input_wrapper["value_info"]

    time_begin = time.time()
    prompt = system_prompt + "\n" + "Task: " + prompt + "\nList of steps: "

    answer = [{'wrapper_id': 'prompt', 'value_type': 'string', 'value_info': prompt}]
    answer = bytes(json.dumps(answer), "utf-8")
    print("Sending data: ", answer)
    channel.basic_publish(exchange='', routing_key=node_name + '_response', body=answer)
