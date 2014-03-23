
# cell -> websocket
sendWS = (x, y, s, over) ->
	buff =  new Uint8Array(3)
	#dv = new DataView(buff)
	buff[0] = x
	buff[1] = y
	ss = 0
	if s is true
		ss = 1

	if over?
		buff[2] = over
	else 
		buff[2] = ss
	
	conn.send(buff)


# cell class with state and init
oncolor = '#FFFF00'
offcolor = '#EEEEEE'
class cell
	@mousedown: false
	@draw: false
	@list: []
	constructor: (@draw, @x, @y, width, border) ->
		@r=draw.rect(width, width).move @x*(width+border), @y*(width+border)
		@off()
		@r.on 'click', @click
		#@r.on 'mousedown', @mdown
		#@r.on 'mouseup', @mup
		@r.on 'mouseenter', @menter
		cell.list.push @
	update: () =>
		if @state == false
			col = offcolor
		else
			col = oncolor
		
		@r.attr {fill: col }
	off: () =>
		@state = false
		@update()
	on: () =>
		@state = true
		@update()
	flipon: () =>
		if not @state
			@on()
			sendWS(@x, @y, @state, null)
	flipoff: () =>
		if @state
			@off()
			sendWS(@x, @y, @state, null)
	flipstate: ()=>
		if @state == true
			@off()
		else
			@on()
		sendWS(@x, @y, @state, null)
	click: () =>
		cell.mousedown = not cell.mousedown
		if cell.mousedown
			@flipstate()
			cell.draw = @state
	mdown: () =>
		cell.mousedown = true
		@flipstate()
	mup: () =>
		cell.mousedown = false
	menter: () =>
		if cell.mousedown is true
			if cell.draw and (not @state)
				@flipon()
			else if not cell.draw and @state
				@flipoff()
	@clear: (y) =>
		if y
			c.on() for c in cell.list
		else
			c.off() for c in cell.list

			

class flipws
	constructor: () ->
		url = "ws://#{location.host}/flipdot"
		console.log("Connecting to " + url)
		@c = new WebSocket(url)
		@c.onopen = @open
	open: () =>
		@status = "open"
		console.log "opened connection to server"
	error: (err) =>
		console.log "error occured #{err}"
	close: () =>
		@status = "closed"
		console.log "connection closed"
	onmessage: (e) =>
		cosole.log "got a message #{e.data}"
	send: (bf) =>
		if @status == "open"
			@c.send(bf.buffer)
			#@c.send("I like pie\n")


conn = new flipws()

$ () ->
	w = 10
	b = 2
	draw = SVG('mainthing').size (w+b)*28, (w+b)*24
	$('#blackbtn').click (e) ->
		cell.clear(false)
		sendWS 0,0,0,224
	$('#yellowbtn').click (e) ->
		cell.clear(true)
		sendWS 0,0,0,240

	new cell(draw, Math.floor(t/24), t%24, w, b) for t in [0..((28*24)-1)]

