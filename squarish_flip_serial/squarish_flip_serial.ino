#define dmask 0x7C
#define bmask 0x1F
#define cmask 0x3F

#define dirpin 13

#define state1 0x5
#define state2 0x6

#define dataDelayMicro 10
#define etime 200

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

void _f(int yellow, int panel) {
  delayMicroseconds(dataDelayMicro);
  PORTC |= (0x4 << panel) | (0x1 << (~yellow & 0x1));
  delayMicroseconds(etime);
  PORTC &= ~cmask;
}

void flip(int x, int y, int yellow, int panel) {
  _x(x);
  _y(y);
  _d(yellow);
  _f(yellow, panel);
}

void clear_dots(int yellow) {
  _d(yellow);
  for(int panel = 0; panel < 3; panel++) {
    for(int y = 0; y < 24; y++) {
      _y(y);
      for(int x = 0; x < 28; x++) {
        _x(x);
        _f(yellow, panel);
        delayMicroseconds(first_flip_delay);
      }
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
  Serial.begin(57600);
}

#define bufsize 512
char buffer[3*bufsize] = {};
int avail, n;
void loop() {
  if((avail = Serial.available()) >= 3) {
    n = 3*min((int)(avail/3), bufsize);
    Serial.readBytes(buffer, n);
    for(int i = 0; i < n; i += 3) {
      switch((int)buffer[i + 2]) {
      case 0xF0:
        clear_dots(1);
        break;
      case 0xE0:
        clear_dots(0);
        break;
      case 0xD0: // Ack request
        Serial.print("H");
        break;
      default:
        flip(buffer[i + 0], buffer[i + 1], buffer[i + 2] & 0x1, buffer[i + 2] >> 1);
      }
    }
  }
}


