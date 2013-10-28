#define dmask 0x7C
#define bmask 0x1F
#define cmask 0x07

#define dirpin 13

#define state1 0x5
#define state2 0x6

#define dataDelayMicro 20
#define etime 400

#define first_flip_delay 300

void resetPins() {
  PORTD &= ~dmask;
  PORTB &= ~bmask;
}

int lookup[28] = {
  0x0 | 0x1,
  0x0 | 0x2,
  0x0 | 0x3,
  0x0 | 0x4,
  0x0 | 0x5,
  0x0 | 0x6,
  0x0 | 0x7,
  0x8 | 0x1,
  0x8 | 0x2,
  0x8 | 0x3,
  0x8 | 0x4,
  0x8 | 0x5,
  0x8 | 0x6,
  0x8 | 0x7,
  0x10 | 0x1,
  0x10 | 0x2,
  0x10 | 0x3,
  0x10 | 0x4,
  0x10 | 0x5,
  0x10 | 0x6,
  0x10 | 0x7,
  0x18 | 0x1,
  0x18 | 0x2,
  0x18 | 0x3,
  0x18 | 0x4,
  0x18 | 0x5,
  0x18 | 0x6,
  0x18 | 0x7
};

void _y(int y) {
  PORTB = lookup[y] | (PORTB & ~bmask);
}

void _x(int x) {
  PORTD = (lookup[x] << 2) | (PORTD & ~dmask);
}

void _d(int yellow) {
  PORTB = ((yellow) ? 0x20 : 0x00) | (PORTB & ~0x20);
  //digitalWrite(dirpin, yellow);
}

void _f(int yellow) {
  delayMicroseconds(dataDelayMicro);
  PORTC |= 0x4 | (0x1 << (~yellow & 0x1));
  delayMicroseconds(etime);
  PORTC &= ~cmask;
}

void flip(int x, int y, int yellow) {
  _x(x);
  _y(y);
  _d(yellow);
  _f(yellow);
}

void clear_dots(int yellow) {
  _d(yellow);
  for(int y = 0; y < 24; y++) {
    _y(y);
    for(int x = 0; x < 28; x++) {
      _x(x);
      _f(yellow);
      delayMicroseconds(first_flip_delay);
    }
  }
}

void serialEvent() {
  if(Serial.available() >= 3) {
    int buffer[3] = {};
    buffer[0] = (int)Serial.read();
    buffer[1] = (int)Serial.read();
    buffer[2] = (int)Serial.read();
    if(buffer[2] == 0x0F) {
      clear_dots(1);
    } else if(buffer[2] == 0x0E) {
      clear_dots(0);
    } else {
      flip(buffer[0], buffer[1], buffer[2]);
    }
  }
}

void setup() {
  DDRD |= dmask;
  DDRB |= bmask;
  DDRC |= cmask;
  PORTD &= ~dmask;
  PORTB &= ~bmask;
  PORTC &= ~cmask;
  pinMode(dirpin, OUTPUT);
  digitalWrite(dirpin, LOW);
  clear_dots(0);
  Serial.begin(9600);
}

void loop() {
  ;
}

