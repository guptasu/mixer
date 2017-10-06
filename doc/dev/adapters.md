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
