
void setup()
{
	Serial.begin(115200);
	while (!Serial)
	{
		; // wait for serial port to connect. Needed for native USB port only
	}
	// pinMode();
	// NOLINT(*)
	Serial1.begin(9600);
	// NOLINTEND
	while (!Serial1)
	{
		;
	}

	/*
	  AT+PARAMETER=<Spreading Factor>,
	  <Bandwidth>,<Coding Rate>,
	  <Programmed Preamble>
	  <Spreading Factor>7~12, (default 12)
	  <Bandwidth>0~9 list as below
	  0 : 7.8KHz (not recommended, over spec.)
	  1 : 10.4KHz (not recommended, over spec.)
	  2 : 15.6KHz
	  3 : 20.8 KHz
	  4 : 31.25 KHz
	  5 : 41.7 KHz
	  6 : 62.5 KHz
	  7 : 125 KHz (default).
	  8 : 250 KHz
	  9 : 500 KHz
	  <Coding Rat>1~4, (default 1)
	  <Programmed Preamble> 4~7(default 4)
	*/
	Serial1.print("AT+PARAMETER=12,7,1,7\r\n");
	delay(1000); // wait for module to respond

	Serial1.print("AT+BAND=915000000\r\n"); // Bandwidth set to 868.5MHz
	delay(1000);							// wait for module to respond

	Serial1.print("AT+ADDRESS=2\r\n"); // needs to be unique
	delay(1000);					   // wait for module to respond

	Serial1.print("AT+NETWORKID=6\r\n"); // needs to be same for receiver and transmitter
	delay(1000);						 // wait for module to respond
	Serial1.print("AT+CRFOP=15\r\n");

	Serial.println("finished setup()");
}

String data;
void loop()
{

	// this set of if statements forwards the uart serial data to the serial monitor and vice versa
	// if (Serial.available() > 0) {      // If anything comes in Serial (USB),
	//   data = Serial.readString();
	//   Serial.print("from serial monitor: " + data);
	//   Serial1.println(data);   // read it and send it out Serial1 (pins 0 & 1)
	// }

	// if (Serial1.available() > 0) {     // If anything comes in Serial1 (pins 0 & 1)
	//   Serial.write(Serial1.read());   // read it and send it out Serial (USB)
	// }

	Serial1.print("AT+SEND=1,6,beacon\r\n");
	if (Serial1.available() > 0)
	{
		data = Serial1.readString();
		Serial.print("from module (send beacon response): " + data);
	}
	delay(3 * 1000);
}
