from dotenv import load_dotenv
import os
import json
import time 
import boto3
import psycopg
import requests
from dataclasses import dataclass, field
from typing import List, Optional, Dict, Set

# TODO: Double Check Code

load_dotenv(dotenv_path='../.env')

s3 = boto3.client('s3')
bucket = "cardimagestorage"



conn = psycopg.connect(os.environ["INSERT_DATABASE_URL"])

cur = conn.cursor()

cur.execute("""
    CREATE TABLE IF NOT EXISTS all_cards_search (
        source_id INT NOT NULL, 
        name TEXT,
        "desc" TEXT,
        humanReadableCardType TEXT,
        race TEXT, 
        atk INT, 
        def INT, 
        level INT,  
        scale INT, 
        attribute TEXT, 
        archetype TEXT, 
        linkval INT, 
        linkmarkers text[],
        banlistInfo text[]
    );
""")

conn.commit()

@dataclass
class Card: 
    id: int = 0
    name: str = "None"
    desc: str = "None"
    pend_desc: str = "None"
    monster_desc: str = "None"
    humanReadableCardType: str = "None"
    race: str = "None"
    atk: int = 0
    defe: int = 0
    level: int = 0
    scale: int = 0
    attribute: str = "None"
    archetype: str = "None"
    linkval: int = 0
    linkmarkers: List[str] = field(default_factory=list)
    banlist_info: Optional[Dict[str, str]] = None

    #not using these
    type: str = "None"
    typeline:List[str] = field(default_factory=list) 
    frameType: str = "None"
    ygoprodeck_url: str = "None"
    card_sets: List[dict] = field(default_factory=list)
    card_images: List[dict] = field(default_factory=list)
    card_prices: List[dict] = field(default_factory=list)
 


    def insertDefaultMonsters(self) -> None: 
        banlist_list = [f"{k}:{v}" for k, v in (self.banlist_info or {}).items()]
        cur.execute("""
            INSERT INTO all_cards_search (source_id, name, humanReadableCardType, "desc", race,level, atk, def, scale, attribute, archetype, linkval, linkmarkers, banlistInfo)
            VALUES (%s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s);
        """, (self.id, self.name, self.humanReadableCardType, self.desc, self.race, self.level, self.atk, self.defe, self.scale, self.attribute, self.archetype, self.linkval, self.linkmarkers, banlist_list))
        conn.commit()



# Get existing source_ids from database
def get_existing_source_ids() -> Set[int]:
    cur.execute("SELECT source_id FROM all_cards_search")
    return {row[0] for row in cur.fetchall()}



existing_ids = get_existing_source_ids()
print(f"Found {len(existing_ids)} existing records in database")

with open("data_new.json", "r") as file: 
    d = json.load(file)
    data = d["data"]

new_cards_count = 0
for card in data[::]:  

    if card["id"] in existing_ids: 
        continue 

    card["defe"] = card.pop('def', 0) 
    c = Card(**card)  
    if "skill" in c.humanReadableCardType.lower(): 
        continue 
    else: 
        c.insertDefaultMonsters()  
        #TODO: Make this a class function 'c.UploadToS3'
        image_url = card["card_images"][0]["image_url"]
        filename = image_url.split("/")[-1]
        response = requests.get(image_url, stream=True)
        if response.status_code == 200:
            # Upload the image directly to S3 without saving locally
            s3.upload_fileobj(response.raw, bucket, filename)
            time.sleep(0.2)
            print(f"Uploaded {filename} to S3")
        else:
            print(f"Failed to download {url}: {response.status_code}") 
        print(card["name"])
        new_cards_count += 1  


print(f"Inserted {new_cards_count} new cards")
