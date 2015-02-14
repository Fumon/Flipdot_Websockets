#include "./font5x7.h"
#include <avr/pgmspace.h>

#define pktsize 4

#define dmask 0x7C
#define bmask 0x1F
#define cmask 0x3F

#define dirpin 13

#define state1 0x5
#define state2 0x6

#define dataDelayMicro 10
#define etime 300

#define first_flip_delay 300

#define pheight 24
#define pwidth 28

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

static unsigned char *fontsByWidth[] = {
  NULL,
  NULL,
  NULL,
  NULL,
  Font5x7
};

int fwidth = 5;
int fheight = 7;
unsigned char *font = &Font5x7[0];

// Cursor (Upper left origin)
int fcx = 0;
int fcy = 0;

void print_character(char c) {
  int x;
  int y;
  unsigned char mask;
  int aindex = ((int)c - 0x20) * fwidth;
  unsigned char line;
  int yellow;
  // TODO: clamp the x and y values so we don't bleed over
  // Loop over character
  for(x = 0; x < fwidth; ++x) {
    line = pgm_read_byte(&Font5x7[aindex + x]);
    for(y = 0, mask = 0x1; y < fheight; ++y, mask = mask << 1) {
      if(line & mask) {
        yellow = 1;
      } 
      else {
        yellow = 0;
      }
      flip((x + fcx) % pwidth, 
      (y + fcy),
      yellow,
      (x + fcx) / pwidth);
    }
  }

  // Change cursor
  fcx += fwidth + 1;
}

void move_cursor(int x, int y) {
  fcx = x;
  fcy = y;
}

void reset_cursor() {
  move_cursor(0,0); 
}

void set_font_dim(int w, int h) {
  font = fontsByWidth[w];
  fwidth = w;
  fheight = h;
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

#define bufsize 256
char buffer[pktsize*bufsize] = {
};
int avail, n;
void loop() {
  avail = Serial.available();
  if(avail > 0 && (avail % pktsize == 0)) {
    n = pktsize*min((int)(avail/pktsize), bufsize);
    Serial.readBytes(buffer, n);
    for(int i = 0; i < n; i += pktsize) {
      switch((int)buffer[i + (pktsize - 1)]) {
      case 0xF0:
        clear_dots(1);
        break;
      case 0xE0:
        clear_dots(0);
        break;
      case 0xD0: // Ack request
        Serial.print("H");
        break;
      case 0xC0: // Print character
        print_character(buffer[i + 0]);
        break;
      case 0xB0: // Move character cursor
        move_cursor((int)buffer[i + 0] + (pwidth * (int)buffer[i + 2]), buffer[i + 1]);
        break;
      case 0xA0: // Reset character cursor
        reset_cursor();
        break;
      case 0x90: // Set font dimentions
        set_font_dim(buffer[i + 0], buffer[i + 1]);
        break;
      case 0x80:
        flip(buffer[i + 0], buffer[i + 1], buffer[i + 2] & 0x1, buffer[i + 2] >> 1);
        break;
      }
    }
  }
}



