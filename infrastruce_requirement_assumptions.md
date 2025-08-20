### Overview
We aim to build an SMS service capable of sending 100 million SMS per month (~38.6 messages/second on average, with potential bursts).

## 1. Time-Based Distribution Calculations

### Monthly Sms Delivery
```
Base: 100,000,000 messages/month

Monthly: 100,000,000 messages
Weekly: 100,000,000 ÷ 4.33 = 23,094,688 messages/week
Daily: 100,000,000 ÷ 30 = 3,333,333 messages/day
Hourly: 3,333,333 ÷ 24 = 138,889 messages/hour
Per Minute: 138,889 ÷ 60 = 2,315 messages/minute
Per Second: 2,315 ÷ 60 = 38.58 messages/second
```

## 2. Golang Web Service RAM Requirements

### Memory Calculation for Go Application

#### Base Application Memory
```
Go Runtime: 50 MB
HTTP Server (Gin/Echo): 10 MB
Database Connections (Pool): 20 MB
Logging & Monitoring: 10 MB
Base Application: ~90 MB
```

#### Request Processing Memory
```
Concurrent Requests Handling:
- Peak concurrent requests: 200 (assuming 5x burst capacity)
- Memory per request: 2 KB (JSON parsing, validation)
- Request memory: 200 × 2 KB = 400 KB

HTTP Response Buffers:
- Response buffer per request: 1 KB
- Total response buffers: 200 × 1 KB = 200 KB
```

#### Message Queue Integration
```
RabbitMQ Client Libraries:
- AMQP client: 5 MB
- Connection pools: 10 MB
- Publishing buffers: 20 MB
- Total queue integration: 35 MB
```

#### Database Integration
```
PostgreSQL Driver (pgx/gorm):
- Driver and ORM: 15 MB
- Connection pool (20 connections): 40 MB
- Query cache: 10 MB
- Total database integration: 65 MB
```

#### Caching & Session Management
```
Redis Client:
- Redis driver: 5 MB
- Connection pool: 5 MB
- Local cache: 20 MB
- Total caching: 30 MB
```

#### Additional Components
```
Rate Limiting Logic: 10 MB
User Authentication Cache: 20 MB
SMS Provider Clients: 15 MB
Monitoring & Metrics: 10 MB
Error Handling & Retry Logic: 5 MB
Total Additional: 60 MB
```

### **Total Go Application RAM Required**
```
Base Application: 90 MB
Request Processing: 1 MB
Queue Integration: 35 MB
Database Integration: 65 MB
Caching: 30 MB
Additional Components: 60 MB
Buffer/Safety Margin: 70 MB

Total per Go Instance: 350 MB
Recommended RAM per instance: 1 GB (3x safety margin)
```

## 3. Redis RAM Requirements

### Redis Memory Calculation

#### User Session Data
```
Active Users: 10,000 businesses
Session data per user: 1 KB (API key, rate limits, last access)
User sessions: 10,000 × 1 KB = 10 MB
```

#### Rate Limiting Counters
```
Rate limit keys per user:
- Daily counter: 1 key × 10,000 users = 10,000 keys
- Hourly counter: 1 key × 10,000 users = 10,000 keys
- Per-minute counter: 1 key × 10,000 users = 10,000 keys

Memory per counter: 50 bytes (key + value + expire)
Total counters: 30,000 × 50 bytes = 1.5 MB
```

#### Message Queue Temporary Storage
```
Temporary message data (before RabbitMQ):
- Buffer size: 5,000 messages
- Data per message: 2 KB
- Queue buffer: 5,000 × 2 KB = 10 MB
```

#### Caching Layer
```
User balance cache: 10,000 × 100 bytes = 1 MB
Provider status cache: 100 × 1 KB = 100 KB
Configuration cache: 5 MB
API response cache: 50 MB
Total caching: 56 MB
```

#### Redis Overhead & Safety
```
Redis base overhead: 50 MB
Replication overhead (if master-slave): 30 MB
Fragmentation buffer: 20 MB
Safety margin: 50 MB
Total overhead: 150 MB
```

### **Total Redis RAM Required**
```
User sessions: 10 MB
Rate limiting: 1.5 MB
Queue buffer: 10 MB
Caching: 56 MB
Redis overhead: 150 MB

Total Redis: 227.5 MB
Recommended RAM: 1 GB (4x safety margin for Redis efficiency)
```

## 4. RabbitMQ RAM Requirements

### RabbitMQ Memory Calculation

#### Base RabbitMQ Memory
```
Erlang VM: 100 MB
RabbitMQ Core: 100-110 MB
Management Plugin: 40 MB
Base system: 180 MB
```

#### Queue Memory Usage
```
Queue Structure:
- sms_normal (70% traffic): 27 msg/sec
- sms_express (30% traffic): 11.5 msg/sec  
- sms_retry (5% of total): 1.9 msg/sec
- sms_failed (1% of total): 0.4 msg/sec

Target queue depth: 5,000 messages per queue (2-minute buffer)
Message size in queue: 2 KB per message

Memory per queue:
- Normal queue: 5,000 × 2 KB = 10 MB
- Express queue: 5,000 × 2 KB = 10 MB
- Retry queue: 2,000 × 2 KB = 4 MB
- Failed queue: 1,000 × 2 KB = 2 MB
- Total queue memory: 26 MB
```

#### Connection & Channel Memory
```
Connections:
- Publishers: 5 connections × 1 MB = 5 MB
- Consumers: 10 connections × 1 MB = 10 MB
- Management: 2 connections × 1 MB = 2 MB

Channels:
- Publisher channels: 20 × 100 KB = 2 MB
- Consumer channels: 30 × 100 KB = 3 MB
- Total connections: 22 MB
```

#### Exchange & Routing Memory
```
Exchanges: 5 × 1 MB = 5 MB
Routing tables: 10 MB
Binding memory: 5 MB
Total routing: 20 MB
```

#### Clustering & Replication (if HA setup)
```
Cluster membership: 10 MB
Queue mirroring: 15 MB (50% of queue memory)
Node coordination: 5 MB
Total clustering: 30 MB
```

### **Total RabbitMQ RAM Required**
```
Base system: 180 MB
Queue memory: 26 MB
Connections: 22 MB
Routing: 20 MB
Clustering: 30 MB
Safety buffer: 70 MB

Total per RabbitMQ node: 348 MB
Recommended RAM per node: 1 GB (3x safety margin)
```

## 5. Complete Infrastructure Summary

### Individual Component RAM Requirements
```
Go Web Service: 1 GB RAM per instance
Redis: 1 GB RAM total
RabbitMQ: 1 GB RAM per node

Minimum Setup:
- 2 Go instances: 2 GB
- 1 Redis instance: 1 GB  
- 1 RabbitMQ node: 1 GB
- Total: 4 GB RAM

Recommended Production Setup:
- 3 Go instances: 3 GB
- 1 Redis master + 1 slave: 2 GB
- 2 RabbitMQ nodes: 2 GB
- Total: 7 GB RAM
```

### Server Configuration Examples

#### Option 1: Single Server (Development/Small Scale)
```
Server Specs:
- RAM: 8 GB
- CPU: 4-8 cores
- Storage: 100 GB SSD

Services:
- 2 × Go instances: 2 GB
- 1 × Redis: 1 GB
- 1 × RabbitMQ: 1 GB  
- OS + Database: 3 GB
- Free RAM: 1 GB
```

#### Option 2: Distributed Setup (Production)
```
Application Server (2 instances):
- RAM: 4 GB each
- CPU: 4 cores each
- Services: 2 × Go web service

Cache Server:
- RAM: 4 GB
- CPU: 2 cores  
- Services: Redis master-slave

Queue Server:
- RAM: 4 GB
- CPU: 4 cores
- Services: RabbitMQ cluster (2 nodes)
```

## 6. Traffic Handling Verification

### At Peak Load (38.58 msg/sec)
```
Go Application:
- Can handle 200 concurrent requests
- Each request processes in ~50ms
- Theoretical capacity: 200 ÷ 0.05 = 4,000 req/sec
- Utilization: 38.58 ÷ 4,000 = 0.96% ✓

Redis:
- Can handle 100,000+ operations/sec
- Our requirement: ~77 ops/sec (2× msg rate)
- Utilization: < 0.1% ✓

RabbitMQ:
- Can handle 10,000+ msg/sec per node
- Our requirement: 38.58 msg/sec
- Utilization: 0.38% ✓
```

### Growth Capacity
```
Current setup can handle up to:
- 10× traffic (1,000M messages/month)
- 385 messages/second
- Before requiring additional servers
```

## 7. PostgreSQL Hard Storage Requirements

### Database Storage Calculation

#### Messages Table (Primary Storage Consumer)
```
Per Message Record Storage:
- Row data: ~671 bytes (as calculated earlier)
- PostgreSQL page overhead: ~24 bytes per row
- Index overhead: ~50 bytes per row (multiple indexes)
- Total per message: ~745 bytes per record

Monthly Storage Growth:
- 100,000,000 messages × 745 bytes = 74.5 GB/month
- Daily growth: 74.5 GB ÷ 30 = 2.48 GB/day
```

#### Index Storage Requirements
```
Primary Indexes (Critical for Performance):
- Primary key (id): ~400 MB/100M records
- user_id index: ~400 MB/100M records  
- created_at index: ~400 MB/100M records
- status index: ~300 MB/100M records
- phone_number index: ~800 MB/100M records
- provider_id index: ~400 MB/100M records

Composite Indexes:
- (user_id, created_at): ~800 MB/100M records
- (status, priority): ~600 MB/100M records
- (provider_id, status): ~600 MB/100M records

Total Index Storage: ~4.7 GB per 100M messages
Index Growth: 4.7 GB/month
```

#### User Table Storage Analysis
```
User Model Fields & Storage:
- ID (uint): 4 bytes
- APIKey (string, 64 chars): 64 bytes
- Name (string, 255 chars): 255 bytes (average 50 chars = 50 bytes)
- Email (string, 255 chars): 255 bytes (average 30 chars = 30 bytes)
- Balance (decimal 12,4): 8 bytes
- RatePerSMS (decimal 6,4): 8 bytes
- DailyLimit (int): 4 bytes
- MonthlyLimit (int): 4 bytes
- IsActive (bool): 1 byte
- CreatedAt, UpdatedAt, DeletedAt (gorm.Model): 24 bytes
- PostgreSQL row overhead: ~24 bytes

Total per user record: ~426 bytes

User Table Storage:
- 10,000 active users × 426 bytes = 4.26 MB
- User growth: ~100 new users/month × 426 bytes = 42.6 KB/month

User Table Indexes:
- Primary key (id): ~40 KB
- Unique index (api_key): ~640 KB (64 bytes × 10K)
- Email index: ~300 KB
- IsActive index: ~10 KB
- Total user indexes: ~990 KB

Total User Table: 4.26 MB + 990 KB = 5.25 MB (minimal growth)
```

#### Provider Table Storage Analysis
```
Provider Model Fields & Storage:
- ID (uint): 4 bytes
- Name (string, 100 chars): 100 bytes (average 20 chars = 20 bytes)
- APIUrl (string, 500 chars): 500 bytes (average 100 chars = 100 bytes)
- APIKey (string, 255 chars): 255 bytes (average 64 chars = 64 bytes)
- APISecret (string, 255 chars): 255 bytes (average 64 chars = 64 bytes)
- CostPerSMS (decimal 6,4): 8 bytes
- Priority (int): 4 bytes
- IsActive (bool): 1 byte
- SuccessRate (decimal 5,2): 8 bytes
- AvgDeliveryTime (int): 4 bytes
- DailyLimit (int): 4 bytes
- CurrentDailyUsage (int): 4 bytes
- LastResetDate (date): 8 bytes
- CreatedAt, UpdatedAt, DeletedAt: 24 bytes
- PostgreSQL row overhead: ~24 bytes

Total per provider record: ~463 bytes

Provider Table Storage:
- 50 SMS providers × 463 bytes = 23.15 KB
- Provider growth: ~2 new providers/year × 463 bytes = 926 bytes/year

Provider Table Indexes:
- Primary key (id): ~200 bytes
- Name index: ~1 KB
- IsActive index: ~50 bytes
- Priority index: ~200 bytes
- Total provider indexes: ~1.45 KB

Total Provider Table: 23.15 KB + 1.45 KB = 24.6 KB (negligible)
```

#### BillingHistory Table Storage Analysis
```
BillingHistory Model Fields & Storage:
- ID (uint): 4 bytes
- UserID (uint): 4 bytes
- TransactionType (string, 20 chars): 20 bytes (average 6 chars = 6 bytes)
- Amount (decimal 12,4): 8 bytes
- BalanceBefore (decimal 12,4): 8 bytes
- BalanceAfter (decimal 12,4): 8 bytes
- Description (text): Variable (average 100 chars = 100 bytes)
- PaymentMethod (string, 50 chars): 50 bytes (average 10 chars = 10 bytes)
- ReferenceID (string, 255 chars): 255 bytes (average 32 chars = 32 bytes)
- MessageID (uint): 4 bytes
- CreatedAt, UpdatedAt, DeletedAt: 24 bytes
- PostgreSQL row overhead: ~24 bytes

Total per billing record: ~232 bytes

BillingHistory Transaction Volume:
With 100M messages/month and billing patterns:
- Credit transactions: ~5,000 users × 2 top-ups/month = 10,000 transactions
- Debit transactions: 100M messages (1 per message) = 100,000,000 transactions
- Refund transactions: ~0.1% of messages = 100,000 transactions
- Total monthly transactions: ~100,110,000 transactions

Monthly BillingHistory Storage:
- 100,110,000 transactions × 232 bytes = 23.23 GB/month
- This is MASSIVE! Much larger than messages table

BillingHistory Indexes:
- Primary key (id): ~4 GB per 100M records
- UserID index: ~400 MB per 100M records
- CreatedAt index: ~800 MB per 100M records
- TransactionType index: ~300 MB per 100M records
- MessageID index: ~400 MB per 100M records
- Total billing indexes: ~5.9 GB/month

Total BillingHistory Growth: 23.23 GB + 5.9 GB = 29.13 GB/month
```

#### PostgreSQL System Storage
```
WAL (Write-Ahead Logging):
- WAL segment size: 16 MB each
- With 100M inserts/month: ~50 WAL files active
- WAL storage: 50 × 16 MB = 800 MB
- WAL archive (if enabled): ~10 GB/month

System Tables & Catalogs:
- pg_stat tables: ~1 GB (statistics)
- System catalogs: ~500 MB
- Temporary files: ~2 GB (for large queries)
- Total system: ~4 GB
```

### **Total PostgreSQL Storage Requirements**

### **Total PostgreSQL Storage Requirements (CORRECTED)**

#### Monthly Growth Breakdown (With All Tables)
```
Messages table data: 74.5 GB/month
Messages table indexes: 4.7 GB/month  
BillingHistory table data: 23.23 GB/month
BillingHistory indexes: 5.9 GB/month
Users table: 5.25 MB (one-time, minimal growth)
Providers table: 24.6 KB (one-time, minimal growth)
WAL & system overhead: 15 GB/month (increased due to more writes)
Total Monthly Growth: ~123.4 GB/month
```

#### Yearly Storage Projection (CORRECTED)
```
Year 1: 123.4 GB × 12 = 1,481 GB (1.48 TB)
Year 2: 1.48 TB + 1.48 TB = 2.96 TB
Year 3: 2.96 TB + 1.48 TB = 4.44 TB
Year 5: 7.4 TB
Year 10: 14.8 TB
```

#### Critical BillingHistory Optimization Strategies

##### Strategy 1: Separate BillingHistory Database
```
Main Database (Messages + Users + Providers): 79.4 GB/month
BillingHistory Database: 29.13 GB/month

Benefits:
- Better performance isolation
- Independent backup/archiving policies  
- Different storage tiers
```

##### Strategy 2: BillingHistory Partitioning
```
Daily Partitions:
- Each partition: ~970 MB/day
- Keep current month active: 29 GB
- Archive older partitions to slower storage
- Partition pruning for old data
```

##### Strategy 3: Summarized Billing (Recommended)
```
Instead of 1 record per message, use:
- Daily summary per user: 10,000 users × 30 days = 300,000 records/month
- Monthly summary: 10,000 records/month
- Keep detailed records only for disputes/audits

Reduced storage: 300,000 × 232 bytes = 69.6 MB/month (99.7% reduction!)
```

##### Strategy 4: Hybrid Approach
```
Real-time billing: Summary records (69.6 MB/month)
Detailed audit trail: Separate audit database or file storage
Query detailed records only when needed for disputes
```

#### Storage Requirements by Retention Policy

##### Scenario 1: Keep All Data (Permanent)
```
Year 1: 1.1 TB
Year 2: 2.2 TB  
Year 3: 3.3 TB
Year 5: 5.5 TB
Recommended Disk: 10 TB (for 7+ years)
```

##### Scenario 2: Archive After 2 Years
```
Active data: 2.2 TB (last 2 years)
Archive storage: Cold storage/S3
Recommended Disk: 5 TB (with 2x buffer)
```

##### Scenario 3: Archive After 1 Year  
```
Active data: 1.1 TB (last year)
Archive storage: Cold storage/S3
Recommended Disk: 3 TB (with 2.5x buffer)
```

### Database Performance Storage Considerations

#### SSD vs HDD Impact
```
With 38.58 inserts/sec continuous:
- Random writes: ~40 IOPS sustained
- Index updates: ~200 IOPS sustained  
- Read queries: ~100 IOPS sustained
- Total: ~340 IOPS sustained

SSD Requirements:
- Any modern SSD (>1000 IOPS) sufficient
- Recommended: Enterprise SSD with 3000+ IOPS

HDD Limitations:
- Standard HDD: ~150 IOPS (insufficient)
- 15K RPM enterprise HDD: ~200 IOPS (barely sufficient)
- RAID 10 HDDs: Could work but not recommended
```

#### Partitioning Storage Strategy
```
Monthly Partitions (Recommended):
- Each partition: ~89 GB
- Old partitions can be moved to slower storage
- Current month on fast SSD: ~90 GB
- Previous 11 months on regular SSD: ~980 GB
- Older data on HDD/cold storage

Weekly Partitions (High Performance):
- Each partition: ~22 GB  
- Current week on NVMe: ~22 GB
- Current month on SSD: ~90 GB
- Older data on regular storage
```

### **Recommended PostgreSQL Storage Setup (UPDATED)**

#### Production Configuration (5 TB Setup - Now Required!)
```
Primary Database Server:
- Fast SSD (NVMe): 1 TB (current year messages + indexes)
- Regular SSD: 3 TB (BillingHistory current year + previous messages)
- HDD/Cold: 1 TB (archive, backups)
- Total: 5 TB

Recommended: Separate BillingHistory to different database/server
```

#### Optimized Configuration with Billing Separation
```
Main Database Server (Messages/Users/Providers):
- SSD Storage: 2 TB (sufficient for 2+ years)

BillingHistory Server:
- SSD Storage: 2 TB (1.5 years of detailed billing)
- HDD Storage: 3 TB (archived billing data)
- Or implement summarized billing: 500 GB total

Total Infrastructure: 4-7 TB depending on strategy
```

#### Budget Configuration (3 TB Setup)
```
Single Server with Optimized Billing:
- Implement summarized billing strategy
- 3 TB SSD storage total:
  - Messages: 80 GB/month growth
  - Optimized billing: 70 MB/month growth
  - Sufficient for 3+ years operation
```

### Storage Requirements Summary by Table

#### Table Storage Breakdown (Per Month)
```
1. Messages Table: 79.2 GB/month (74.5 GB data + 4.7 GB indexes)
2. BillingHistory Table: 29.13 GB/month (23.23 GB data + 5.9 GB indexes)  
3. Users Table: ~5 MB total (minimal growth)
4. Providers Table: ~25 KB total (minimal growth)
5. System/WAL: 15 GB/month

Total: 123.4 GB/month growth
CRITICAL: BillingHistory is 24% of total storage!
```

### Backup Storage Requirements
```
Full Backup Size (monthly):
- Compressed backup: ~25% of original = 22 GB/month growth
- Keep 12 monthly backups: 22 GB × 12 = 264 GB
- Daily incremental backups: ~2.5 GB × 30 = 75 GB
- Total backup storage needed: 340 GB

Recommended Backup Storage: 500 GB separate disk/location
```
