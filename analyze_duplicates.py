#!/usr/bin/env python3
"""
Quick script to analyze duplicates in the property data
"""

import json

# Sample of the JSON data from the MCP response (first few properties)
sample_data = [
    {"id": "prop_2714hipawaiplacehono", "address": "2714 Hipawai Place", "price": 1600000},
    {"id": "prop_2714hipawaiplacehono", "address": "2714 Hipawai Place", "price": 1600000},
    {"id": "prop_2772kalawaostunit29h", "address": "2772 Kalawao St Unit 29", "price": 1942000},
    {"id": "prop_2772kalawaostunit29h", "address": "2772 Kalawao St Unit 29", "price": 1942000},
    {"id": "prop_3122kaloaluikisthono", "address": "3122 Kaloaluiki St", "price": 1850000},
    {"id": "prop_3122kaloaluikisthono", "address": "3122 Kaloaluiki St", "price": 1850000},
]

# Count duplicates by ID
id_counts = {}
for prop in sample_data:
    prop_id = prop["id"]
    if prop_id in id_counts:
        id_counts[prop_id] += 1
    else:
        id_counts[prop_id] = 1

print("Duplicate Analysis:")
print("==================")
for prop_id, count in id_counts.items():
    if count > 1:
        address = next(p["address"] for p in sample_data if p["id"] == prop_id)
        print(f"ID: {prop_id}")
        print(f"Address: {address}")
        print(f"Count: {count}")
        print("---")

# Based on the MCP response, I can see that the 243 properties contain many duplicates
print("\nFrom the MCP response, I observed:")
print("- Total properties returned: 243")
print("- Pages scraped: 3 (116 + 118 + 9)")
print("- Many properties appear exactly twice")
print("- This suggests pages have overlapping content or parsing issues")
