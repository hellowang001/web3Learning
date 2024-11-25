import schedule
import time


def job():
    print("I m working")


schedule.every().day.at("09:00").do(job)

i = 0
while i < 10:
    schedule.run_pending()
    print("Schedule Running 1 ", i)
    time.sleep(1)
    i = i + 1
