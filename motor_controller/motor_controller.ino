#include <Arduino.h>

// Debug
#define DEBUG           1
#ifndef DEBUG
#define Serial          Serial1
#endif

// Pins
#define LEFT_FORWARD    9
#define LEFT_REVERSE    10
#define RIGHT_FORWARD   11
#define RIGHT_REVERSE   12

// Messaging
#define SPLIT_CHAR              ':'
#define END_CHAR                '!'
#define LEFT_COMMAND            "L"
#define RIGHT_COMMAND           "R"
#define PARSE_ERROR             "PARSE ERR"
#define PARAM_ERROR             "OUT OF RANGE"
#define COMMAND_ERROR           "UNKNOWN COMMAND"
#define BAD_COMMAND             ""
#define BAD_PARAM               -9999
#define SERIAL_BAUD             115200

typedef struct {
    String command;
    int param;
} command_t;

command_t current_command;

void setup() {
    Serial.begin(SERIAL_BAUD);
}

void loop() {
    while (Serial.available()) {
        String message = Serial.readStringUntil(END_CHAR);

        parse_message(&current_command, message);

        if (current_command.command.length() == 0) {
            send_reply(PARSE_ERROR);
            break;
        }

        if (current_command.param < -255 || current_command.param > 255) {
            send_reply(PARAM_ERROR);
            break;
        }

        // Command is good, do something with it
        if (current_command.command == LEFT_COMMAND) {
            if (current_command.param >= 0) {
                analogWrite(LEFT_REVERSE, 0);
                analogWrite(LEFT_FORWARD, current_command.param);
            }
            else {
                analogWrite(LEFT_FORWARD, 0);
                analogWrite(LEFT_REVERSE, current_command.param * -1);
            }

        }
        else if (current_command.command == RIGHT_COMMAND) {
            if (current_command.param >= 0) {
                analogWrite(RIGHT_REVERSE, 0);
                analogWrite(RIGHT_FORWARD, current_command.param);
            }
            else {
                analogWrite(RIGHT_FORWARD, 0);
                analogWrite(RIGHT_REVERSE, current_command.param * -1);
            }
        }
        else {
            send_reply(COMMAND_ERROR);
            break;
        }

        // echo command back to sender on success
        send_reply(message);
    }
}

void parse_message(command_t* com, String message) {
    int split_index = message.indexOf(SPLIT_CHAR);

    if (split_index == -1) {
        com->command = BAD_COMMAND;
        com->param = BAD_PARAM;
    }
    else {
        com->command = message.substring(0, split_index);
        com->param = message.substring(split_index+1).toInt();
    }
}

void send_reply(String message) {
    Serial.print(message);
    Serial.print(END_CHAR);
}