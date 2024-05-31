import time

import pika
import json

from exllamav2 import(
    ExLlamaV2,
    ExLlamaV2Config,
    ExLlamaV2Cache,
    ExLlamaV2Tokenizer,
)

from exllamav2.generator import (
    ExLlamaV2BaseGenerator,
    ExLlamaV2Sampler
)

model_directory = "..\\models\\Llama-2-13B-GPTQ"

config = ExLlamaV2Config(model_directory)
model = ExLlamaV2(config)
cache = ExLlamaV2Cache(model, lazy=True)
model.load_autosplit(cache)
tokenizer = ExLlamaV2Tokenizer(config)

generator = ExLlamaV2BaseGenerator(model, cache, tokenizer)

settings = ExLlamaV2Sampler.Settings()
settings.temperature = 0.85
settings.top_k = 50
settings.top_p = 0.8
settings.token_repetition_penalty = 1.01
settings.disallow_tokens(tokenizer, [tokenizer.eos_token_id])

max_new_tokens = 100
generator.warmup()

connection = pika.BlockingConnection()
channel = connection.channel()

node_name = 'llama2'

channel.queue_declare(queue=node_name)
channel.queue_declare(queue=node_name+'_response')

print("Ready to handle messages")
while True:
    method_frame, properties, body = next(channel.consume(node_name))
    channel.basic_ack(method_frame.delivery_tag)

    data = json.loads(body.decode("utf-8"))
    #print("Request: ", data)

    system_prompt = ""
    prompt = ""
    prompt_list = []
    prompt_steps = []
    list_process_mode = ""

    answer = []
    for input_wrapper in data:
        if input_wrapper['wrapper_id'] == "prompt":
            prompt = input_wrapper["value_info"]
        if input_wrapper['wrapper_id'] == "system_prompt":
            system_prompt = input_wrapper["value_info"]
            answer.append(input_wrapper)
        elif input_wrapper['wrapper_id'] == "prompt_list":
            prompt_list.append(input_wrapper["value_info"])
        elif input_wrapper['wrapper_id'] == "prompt_steps":
            prompt_steps.append(input_wrapper["value_info"])
            answer.append(input_wrapper)
        elif input_wrapper['wrapper_id'] == "max_new_tokens":
            max_new_tokens = int(input_wrapper["value_info"])
        elif input_wrapper['wrapper_id'] == "list_process_mode":
            list_process_mode = input_wrapper["value_info"]
        else:
            answer.append(input_wrapper)

    if prompt_list:
        if list_process_mode == "steps":
            time_begin = time.time()
            pre_prompt = ""
            for idx, elem in enumerate(prompt_steps):
                pre_prompt = pre_prompt + f"Step {idx+1} result: " + elem + "\n"
            prompt = prompt_list[0] + "\nAnswer: "
            if system_prompt:
                prompt = system_prompt + prompt
            prompt = pre_prompt + prompt
            output = generator.generate_simple(prompt, settings, max_new_tokens)
            #print("----------------------------")
            #print(output)
            print("----------------------------")
            output = prompt + output[len(prompt):].split(sep="\n")[0]
            print(output)
            print("----------------------------")
            time_end = time.time()
            time_total = time_end - time_begin
            print(f"Response generated in {time_total:.2f} seconds, {max_new_tokens} tokens, {max_new_tokens / time_total:.2f} tokens/second")

            if len(prompt_list) > 1:
                answer.append({'wrapper_id': 'prompt_steps', 'value_type': 'string', 'value_info': output[len(prompt):]})
                for i in range(len(prompt_list)):
                    if i != 0:
                        answer.append({'wrapper_id': 'prompt_list', 'value_type': 'string', 'value_info': prompt_list[i]})
                answer.append({'wrapper_id': 'list_process_mode', 'value_type': 'string', 'value_info': 'steps'})
            else:
                answer.append({'wrapper_id': 'answer', 'value_type': 'string', 'value_info': output[len(prompt):]})
        if list_process_mode == "parallel":
            for elem in prompt_list:
                prompt = elem
                time_begin = time.time()
                output = generator.generate_simple(prompt, settings, max_new_tokens)
                print("----------------------------")
                print(output)
                print("----------------------------")
                time_end = time.time()
                time_total = time_end - time_begin
                print(f"Response generated in {time_total:.2f} seconds, {max_new_tokens} tokens, {max_new_tokens / time_total:.2f} tokens/second")

                answer.append({'wrapper_id': 'answer_list', 'value_type': 'string', 'value_info': output[len(prompt):]})
    else:
        time_begin = time.time()
        output = generator.generate_simple(prompt, settings, max_new_tokens)
        print("----------------------------")
        print(output)
        print("----------------------------")
        time_end = time.time()
        time_total = time_end - time_begin
        print(f"Response generated in {time_total:.2f} seconds, {max_new_tokens} tokens, {max_new_tokens / time_total:.2f} tokens/second")

        answer.append({'wrapper_id': 'answer', 'value_type': 'string', 'value_info': output[len(prompt):]})

    answer = bytes(json.dumps(answer), "utf-8")
    print("Sending data: ", answer)
    channel.basic_publish(exchange='', routing_key=node_name + '_response', body=answer)
