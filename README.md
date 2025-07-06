# Ravant_Media_Assesment

# 1. Concurrency in Golang

This project demonstrates how to process a large number of tasks (e.g., file operations) efficiently using **Golang concurrency patterns** — specifically **Goroutines**, **Channels**, and a **Worker Pool**.

---

## Real-world context

In a production batch job I worked on:

- We **fetched ~100 files at once** from a feed hub.
- For each file:
  - I **read the file**
  - I **validated its contents**
  - I **wrote it to MongoDB**
- I set up a **worker pool** using Goroutines where each file was processed concurrently:
  - A **buffered channel** fed file paths into the pool.
  - **Workers (Goroutines)** picked up file paths and ran the read-validate-write sequence.
- This batch job ran as an **OpenShift CronJob**, scheduled **monthly or quarterly**.

This approach ensured:
- Efficient resource usage without overwhelming the system
- Scalable and predictable performance
- Safe concurrent processing

---