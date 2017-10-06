# Writing Mixer Adapters

This will eventually turn into a developer's guide for 
creating Mixer adapters. For now, it's just a set of
notes and reminders:

![test](https://github.com/guptasu/mixer/blob/ADG/doc/dev/img/Template%2C%20Adapter%20and%20Operator%20config%20relationship.svg)
![test](https://github.com/guptasu/mixer/blob/ADG/doc/dev/img/Mixer%20Start%20Flow.svg)
![test](https://github.com/guptasu/mixer/blob/ADG/doc/dev/img/Operator%20config%20change.svg)
![test](https://github.com/guptasu/mixer/blob/ADG/doc/dev/img/Request%20time%20.svg)

- Adapters must use env.Logger for logging during
execution. This logger understands about which adapter
is running and routes the data to the place where the
operator wants to see it.

- Adapters must use env.ScheduleWork or env.ScheduleDaemon
in order to dispatch goroutines. This ensures all adapter goroutines
are prevented from crashing Mixer as a whole by catching
any panics they produce.


# Summary diagrams

Below diagrams show the relationship between, Adapter, Mixer, Template and operator config. It also shows the flow of Mixer at boot time, how it interacts with Adapters and operator configuration. The diagrams also demonstrate how handler, rule and instance configuration of operator config is translated to calls into Adapters during config load time and request time.

First let's look into how Mixer, adapter, templates and operator configurations are related

![template, adapter and operator config relationship](https://github.com/guptasu/mixer/blob/ADG/doc/dev/img/Template%2C%20Adapter%20and%20Operator%20config%20relationship.svg)

Now we have understood the relationship between various artifacts, let's look into what happens at the time of Mixer start, when operator configuration is loaded/changed and when request is received.

![flow: mixer start](https://github.com/guptasu/mixer/blob/ADG/doc/dev/img/Mixer%20Start%20Flow.svg)

![flow: operator config change](https://github.com/guptasu/mixer/blob/ADG/doc/dev/img/Operator%20config%20change.svg)

![flow: incoming api](https://github.com/guptasu/mixer/blob/ADG/doc/dev/img/Request%20time%20.svg)

