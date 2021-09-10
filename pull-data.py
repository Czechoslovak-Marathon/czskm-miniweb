import requests
import csv

themes = [
    'Retro (<2000)', 'Horror', 'České hry', 'Racing hry',
    'Blind runy', 'GTA-like runy', 'Meme runy', 'Co-op runy',
    'Indie hry', 'Hry s českým dabingem', 'Filmové hry',
    'Disney, Pixar, apod.', 'Bad hry', 'Platformer hry',
    '100% runy', 'Modern hry (>2010)', 'PC hry',
    'Console hry', '<30 min runy', 'Glitchless runy',
    'GoldSrc / Source hry', 'PS1 hry', 'Remaky'
]

csv_file = requests.get('https://docs.google.com/spreadsheets/d/1CkXButhg9CzPfvfAUWRfKOMvvmb1ZA8EbX3g1IElaKw/gviz/tq?tqx=out:csv')
entries = [x.split(',') for x in csv_file.text.splitlines()[1:]]

ratings = {}

i = 0
print(f'Number of entires: {len(entries)}')
for theme in themes:
    rating = 0
    runs = 0
    for entry in entries:
        rating += int(entry[i+1].replace('"', ''))
        if entry[i+2].strip() == '"Ano"':
            runs += 1
    ratings[theme] = (rating, runs)
    i += 2

sort_themes = sorted(ratings.items(), key=lambda x: x[1][0]+3*x[1][1], reverse=True)

with open('/root/czskm-miniweb/results.csv', 'w', encoding='utf-8', newline='') as f:
    writer = csv.writer(f)
    for theme in sort_themes:
        writer.writerow([theme[0],theme[1][0],theme[1][1],theme[1][0]+3*theme[1][1]])