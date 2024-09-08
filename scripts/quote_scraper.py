import requests, re, sys, json, time
from bs4 import BeautifulSoup as bs


HEADERS = {'User-Agent': 'Mozilla/4.0 (iPad; CPU OS 12_1 like Mac OS X) AppleWebKit/603.1.15 (KHTML, like Gecko) Mobile/15E148'}
class QuoteScraper:
    def __init__(self, max_tags_per_quote=10, output_file=sys.stdout):
        self.output_file = output_file
        self.max_tags = max_tags_per_quote
        self.quote_idx = 0
        self.page_idx = 0
        self.current_page_quotes = self.get_current_quotes()
        if self.current_page_quotes is None:
            raise "Quote scraper initialisation failed"


    def next_quote(self):
        quote_author_text = self.current_page_quotes[self.quote_idx].find_next("span").extract().get_text(strip=True)
        source_title = ""
        source_type = ""
        if ',' in quote_author_text:
            author = quote_author_text.strip(',')
            source_title = self.current_page_quotes[self.quote_idx].find_next("span").extract().get_text(strip=True)
            source_type = "Book"
        else:
            author = quote_author_text
        
        tags = [a.text for a in self.current_page_quotes[self.quote_idx].find_next("div", class_="greyText smallText left").find_all("a")]
        
        content = self.current_page_quotes[self.quote_idx].get_text(separator='\n', strip=True)
        stripped_content = re.search('“(.*?)”', content, re.DOTALL).group(1)

        self.quote_idx += 1
        if self.quote_idx >= len(self.current_page_quotes):
            self.page_idx += 1
            time.sleep(0.5) # limit request rate to not get ip banned
            self.current_page_quotes = self.get_current_quotes()
            self.quote_idx = 0

        quote_dict = {
            "content": stripped_content,
            "author": author,
            "source": {"title": source_title, "type": source_type},
            "tags": tags[0:self.max_tags],
        }
        return quote_dict

    def output_next_quote(self, output_file=None):
        if output_file is None:
            output_file = self.output_file
        next_quote = self.next_quote()
        json.dump(next_quote, output_file)
        print("", file=output_file)

    def get_current_quotes(self):
        print(f"Getting page {self.page_idx} of quotes")
        response = requests.get(f"https://www.goodreads.com/quotes?page={self.page_idx}", headers=HEADERS)
        return bs(response.content, "html.parser").find_all("div", class_="quoteText")
    
if __name__ == "__main__":
    num_quotes = int(input("Please enter the number of quotes to scrape: "))
    with open('quotes.json', 'w') as f:

        scraper = QuoteScraper(output_file=f)
        for i in range(num_quotes):
            scraper.output_next_quote()
        
        f.close()