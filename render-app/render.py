import bpy

bpy.context.scene.render.filepath = "E:/projects/render-service/render-app/output.png"
bpy.context.scene.render.image_settings.file_format = 'PNG'
bpy.ops.render.render(animation=False, write_still=True)
