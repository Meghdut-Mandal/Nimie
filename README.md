

# Nimie


## An End-to-end encrypted Anonymous Messaging Service.


## In a nutshell

This service basically empowers you to have short anonymous conversations with people, with the added security of end-to-end encryption. The interface is intuitive enough for users to use it like a regular IM app.


## User Story

Lupita, an undergrad student of kiit, is tired of using whatsapp and scrolling through instagram. Most of her conversions over whatsapp are mundane and mostly about regular college stuff. To have some fun and possibly get some honest compliments from people, she decides to try out the Nimie App. Via the app she generates a unique link and shares it over her whatsapp status. Within a few hours, her app inbox gets filled with various messages ranging from “I love you” to “Kaam dhanda nahi hai kiya? ”. She can then tap on the “I love you” text and replies back with a painful “Sorry I’ve a bf”. The anonymous person(aka simp) then gets a notification of the same. Lupita shares a screenshot of the conversation on her whatsapp for flexing. After 24 hours the link and the message box expires.


## Proposed Architecture



* Asymmetric Cryptography( public key and private key encryption)
* Discord Threads as backend (in version 2)
* Normal Golang +sqlite backend (in version 1)
* Redis
* Cloudflare Workers for hosting the frontend.


## Technical Overview


### Glossary :



* User : Unique id + public key (No other information should be taken from the user!)
* Public key: A key that can be used to encrypt data, and is shared across devices
* Private Key: A key that can be used to decrypt data and remains on the user device.
* Status: This is analogous to the **status** feature we have in other IM apps which expire after 24 hours. When the user(aka user A) wants to have an anonymous conversation he can create a status(in the Nimie app) with some text  and share it across other social media. The other user(aka user B) who opens the link can then chat with user A. (Unencrypted text + unique link + user id + creation time)
* Conversation: When user B replies back to user A a conversation begins. Conversations are short lived message threads between 2 users. The max size of a text message in conversion is 500 characters. However they can send multiple texts in a chain.
* This app isn't meant for having long conversations, as a human tendency having long conversations is either useless or in order to make it meaningful we give away enough context information to determine each other's identity. Users should consider switching to something personal like Signal or telegram.
* To test the websockets, you can use the [WebSockets King Client Chrome Extension](https://chrome.google.com/webstore/detail/websocket-king-client/cbcbkhdmedgianpaifchdaddpnmgnknn/related?hl=en).


### The user has the following use cases



1. Open the app, generate keys (A random name and avatar stays on the device though) and register as a User, this won't require any authentication.
2. Register with the server, exchanging public keys.
3. The server assigned a unique client id which the client will use in future to communicate.
4. Create a new status with text.
5. A conversation starts by someone replying back to a status.
6. See recent conversations.
7. Continue recent conversations.


### Objects



1. Conversation : (**conversation_id**, **user_id_a**, **user_id_b**, **created_at**, **status_id**)
2. User (**user_id**, **create_time**, **public_key**)
3. ChatMessage( **message_id**, **conversation_id**, **create_time**, **user_id**, **message**)
4. Status(**status_id**, **create_time**, **header_text**, **user_id**, **link_id**)


### API Overview



* Public register API - Doesn't need any authorization. Rest of the api’s should be protected by expirable JWT token.
* The apis will keep evolving to **reduce** the amount of data being exchanged, be it server to client or client.
* JWT should be used to authenticate - not yet done actually. However, having said that the messages are encrypted and therefore even if a third party or the developer tries they can't decrypt the content.
* The conversations history and stuff don't leave the device, the conversation id are stored on the device and **can't** be retrieved from the backend Api.


### API - Docs


<table>
  <tr>
   <td>Name
   </td>
   <td>Path
   </td>
   <td>Parameters/Body  
   </td>
   <td>Response / Notes 
   </td>
  </tr>
  <tr>
   <td>Register
<p>
New User 
<p>
This doesn't require auth.
<p>
DONE
   </td>
   <td>POST
<p>
/user/register
   </td>
   <td>Json body
<p>
{
<p>
“public_key” : “dcdwc"
<p>
}
<p>
The public key should be base64 encoded PKCS1 Public key. The request will fail otherwise.
   </td>
   <td>
<ul>

<li><strong>201</strong> : Json 
    {
<p>

    “user_id”: 96493264,
<p>

    “create_time” :  232323323,
<p>

    “token_jwt”: “de34e3e3”
<p>

    }
<ul>

<li><strong>400</strong>: Invalid Public key.
</li>
</ul>
</li>
</ul>
   </td>
  </tr>
  <tr>
   <td>Create Status
<p>
DONE
   </td>
   <td>POST:
<p>
/status/create
   </td>
   <td>Json Body:
<p>
{
<p>
“text” : “ Nice status heading”
<p>
}
   </td>
   <td>
<ul>

<li>201: JSON
    {
<p>

    “unique_id”: “sdsdbkbds”
<p>

    }
</li>
</ul>
   </td>
  </tr>
  <tr>
   <td>Delete Status
   </td>
   <td>Delete:
<p>
/status/{<strong>id</strong>}
   </td>
   <td>nil
   </td>
   <td>
<ul>

<li>200
</li>
</ul>
   </td>
  </tr>
  <tr>
   <td>Initiate Conversation
   </td>
   <td>POST:
<p>
/status/reply/{id}
   </td>
   <td>Json Body:
<p>
{
<p>
“text”: “ Nice Status!”,
<p>
 “b_key” :”dwd33e”
<p>
}
   </td>
   <td>
<ul>

<li>201: Json 
   {
<p>
    “conversation_id”: “134bcsdc”,
<p>
    “a_key”: “cw242d…”
<p>
   }
</li>
</ul>
   </td>
  </tr>
  <tr>
   <td>Connect 
<p>
Session WSS
<p>
This url starts with a ws://… 
<p>
Representing a websocket connection.
<p>
DONE
   </td>
   <td>POST:
<p>
/conversation/{<strong>conversation_id</strong>}/{<strong>user_id</strong>}
   </td>
   <td>nil	
<p>
<strong>conversation_id</strong>,<strong>user_id</strong> are long integers.  
   </td>
   <td>
<ul>

<li>Upgrade Request to ws

<li>403 if the user isn't part of the conversation		

</li>
</ul>
   </td>
  </tr>
  <tr>
   <td>End
<p>
Conversation
   </td>
   <td>POST
<p>
/conversation/close/{id}
   </td>
   <td>nil
   </td>
   <td>
<ul>

<li>200 
</li>
</ul>
   </td>
  </tr>
</table>

