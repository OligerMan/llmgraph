import time

import pika
import json

connection = pika.BlockingConnection()
channel = connection.channel()

node_name = 'list_parser'

channel.queue_declare(queue=node_name)
channel.queue_declare(queue=node_name+'_response')

print("Ready to handle messages")
while True:
    method_frame, properties, body = next(channel.consume(node_name))
    channel.basic_ack(method_frame.delivery_tag)

    data = json.loads(body.decode("utf-8"))
    print("Request: ", data)

    answer = []

    for elem in data:
        if elem["wrapper_id"] == "answer":
            parsed_data = elem["value_info"].split(sep="\n\n")[0]  # delete everything after the list
            parsed_data = parsed_data.split(sep="\n")              # split to steps
            for step in parsed_data:
                if step:
                    step = step.split(sep=") ", maxsplit=1)[1]
                    answer.append({'wrapper_id': 'prompt_list', 'value_type': 'string', 'value_info': step})

    answer = bytes(json.dumps(answer), "utf-8")
    print("Sending data: ", answer)
    channel.basic_publish(exchange='', routing_key=node_name + '_response', body=answer)
