# LeakTrack
LeakTrack is a security oriented product designed to locate, remove and notify the user of PII images from various image-sharing sites.


## How does it work
LeakTrack scans in real-time the images uploaded to image-sharing sites
for potentially identifiable information, when located - any users found matching
the identity of the owner will get a notification of the type of document, where it
was found and minimal information to be able to verify the leak and disregard false positives.

## Privacy
Your privacy is our top priority. All images containing PII are heavily redacted and only minimal amount of information is sent 
back to our users for the purpose of false positive elimination.

For example, if we identified a credit card that potentially belongs to you, we will
send you the last 4 digits for you to verify.

Similarly, for passports / id cards, we would send the last 4 digits of
the unique identifier.

Signing up to LeakTrack does not require your data at all, anyone is able to request
a unique tracking number which will be used to view your personal events & set up your triggers.

You may also delete your account at any time, and all previously held triggers from your account
will be deleted with it.


## REST Api
If you would like to integrate LeakTrack into your own projects,
products feel free to use our hosted api.

### Rate Limit
To not exceed our free hosting limits, we had to introduce rate limiting to the API.

-   PUT/POST/DELETE -- 100 packets/s per ip
-   GET -- 500 packets/s per ip


### Terminology
Throughout the application you may see various terms pop up, here's a quick rundown of what we mean by them:

- **Triggers**: User-defined queries to select which data the user gets notified on.
- **Events**: Notifications of discovery of potential leaks.
- **PII**: Personally Identifiable Information. 