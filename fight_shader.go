package main

import "github.com/Noofbiz/pixelshader"

var fightShader = &pixelshader.PixelShader{FragShader: `
  //Modified version of Shader from ShaderToy
  //Modified by: Noofbiz
  //Created by: PetrifiedLasagna
  //shader source: https://www.shadertoy.com/view/Mdc3WH
  #ifdef GL_ES
  #define LOWP lowp
  precision mediump float;
  #else
  #define LOWP
  #endif
  uniform vec2 u_resolution;  // Canvas size (width,height)
  uniform vec2 u_mouse;       // mouse position in screen pixels
  uniform float u_time;       // Time in seconds since load
  uniform sampler2D u_tex0;   // Drawable Tex0
  float flashrad = 75.0;
  float flash_min_dist = 1.5;
  float flash_falloff_start = 3.0;
  float flash_falloff_end = 6.0;
  void main()
  {
    vec2 p = gl_FragCoord.xy/u_resolution.xy;
    vec2 texPos = vec2(p.x, -p.y);
    gl_FragColor = texture2D(u_tex0, texPos);

    if(flashrad == 0.0)
    return;

    vec2 flashp = u_mouse.xy;
    float dist = abs(sin(u_time/4.0))*5.5;
    float light;
    if(flash_falloff_end > 0.0 && dist>flash_falloff_end)
        light = flashrad*flash_falloff_end;
    else if(dist<flash_min_dist)
        light = flashrad*flash_min_dist;
    else
        light = flashrad*dist;

    light *= light;

    float dcol = pow(gl_FragCoord.x - flashp.x, 2.0) + pow(gl_FragCoord.y - flashp.y, 2.0);

    if(dcol >= light *.25){
    	if(dcol < light *.75)
        	gl_FragColor *= .75;
        else if(dcol < light)
            gl_FragColor *= .5;
        else
            gl_FragColor *= .018;
    }

    if(flash_falloff_end > 0.0 && dist >= flash_falloff_start){
        float scalar = 1.0 - clamp((dist - flash_falloff_start) / (flash_falloff_end - flash_falloff_start),
                                 0.0, 1.0);
        gl_FragColor *= scalar;
    }

    gl_FragColor.a = 1.0;
  }
`}
