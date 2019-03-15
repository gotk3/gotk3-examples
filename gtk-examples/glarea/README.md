GLArea example
==============

The glarea example shows how to create a GLArea widget, draw to it, and 
animate the contents inside of a normal gotk3 application.

## Dependencies

Even though GLArea is a widget that is a part of gotk3, it is merely a canvas
that shows the results of OpenGL api calls. Sadly, that means that this example 
will not run out of the box, and additional requirements are necessary. In order
to draw and animate the contents of the GLArea this example uses two wonderful 
libraries from go-gl.

### go-gl/gl

[`gl`][1] enables the necessary OpenGL api calls to get things drawn to the GLArea 
widget, and all from a go context.

### go-gl/mathgl

[`mathgl`][2] makes the math surrounding simpler by providing convenient data 
structures and functions.


#### Final note about the additional dependencies. 

The GLArea widget can be embedded inside of an application with no extra libraries
and it will work just fine. However, their will just be a blank canvas inside the 
application, and although that can be exciting on some level, it would fail to 
really demonstrate the true implications of the widget.


### Screenshot

![glare.png](https://github.com/drakbar/gotk3-examples/blob/master/gtk-examples/glarea/glarea.png)

### Credits
- [dmitshur][3]
- [errcw][4]
- [depy][5]
- [alexlarsson][6]
- [baedert][7]
- [ebassi][8]
- [khronos.org][9]

[1]:https://github.com/go-gl/gl
[2]:https://github.com/go-gl/mathgl
[3]:https://github.com/dmitshur
[4]:https://github.com/errcw
[5]:https://github.com/depy
[6]:https://github.com/alexlarsson
[7]:https://github.com/baedert
[8]:https://github.com/ebassi
[9]:https://www.khronos.org/opengl/wiki/
