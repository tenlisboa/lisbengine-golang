#version 330 core
layout (location = 0) in vec3 aPos;
layout (location = 1) in vec3 aTexCoord;

out vec3 ourColor;
out vec2 TexCoord;

uniform mat4 model;
uniform mat4 view;
uniform mat4 projection;

void main()
{
   TexCoord = vec2(aTexCoord.xy);
   gl_Position = projection * view * model * vec4(aPos, 1.0f);
}