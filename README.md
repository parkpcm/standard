# ParkPCM Golang Standard Modal



### [Database](database)

Provides stuct to connect to mySQL database, using standard credentials GCP Secret Manager finder with via SECRET_PATH environment variable or file mounted volume

Connections methods depending on service rollout:

- Connect via private IP Address
- Connect via Socket

Data connection secret JSON format:

```
 {
  "username":"",
  "password":"",
  "instance":"",
  "database":"",
  "public":"",
  "private":""
} 
```
---

### [Email](email)

Create a MailGun client using standard credentials GCP Secret Manager finder with via SECRET_PATH environment variable or file mounted volume

Email secret JSON format:

```
 {
  "mailgun_key":"",
  "mailgun_domain":""
} 
```

---

### [Secret](secret)

Returns a `[]byte` of the GCP Secret Manager

**Get** - accepts path with context

**GetFromVolume** - accepts secret mounted as file

----