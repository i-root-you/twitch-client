# OBS Websocket Client
Just doing a quick code review I already saw signs of pretty bad porgramming
practices such as short polling, which is definitely not necessary when we are
using something like websockets, we should be using event loops, or some form of
async design model instead of wasting resources.

So I'm just going to pull down the client and modify it directly as I work,
removing short polling and whatever else is bad. Just hope its not so bad that I
should have just wrote it from scratch or just leveraged the C++ version. I have
not yet done a full code review so I can not say for certain.

