#!/usr/bin/python3
from icalendar import Calendar
import json
import sys

file_name = sys.argv[1]

ics = open(file_name, "r")
event_id_to_name = {
    "Segregowane": "odpady segregowane ♻️",
    "biodegradowalne": "odpady BIO 🍂",
    "komunalne": "odpady zmieszane 🚮",
    "gabaryty": "odpady wielogabarytowe",
}

cal = Calendar().from_ical(ics.read())
events = {}
for component in cal.walk():
    if component.name == "VEVENT":
        date_string = component.get("dtstart").dt.strftime('%Y-%m-%d')
        event_type = event_id_to_name[str(component.get("summary"))]
        if not date_string in events:
            events[date_string] = []
        events[date_string].append(event_type)

dump_events = []
for i in events:
    dump_events.append({"date": i, "events": events[i]})
print(json.dumps(dump_events, indent=4, ensure_ascii=False))
