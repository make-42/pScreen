gnome-terminal -- bash -c "openscad -o ./STLs/SSD1322_bottom.stl -p ./OpenSCAD/bottom.json -P SSD1322 ./OpenSCAD/bottom.scad; exit"
gnome-terminal -- bash -c "openscad -o ./STLs/SSD1322_top.stl -p ./OpenSCAD/top.json -P SSD1322 ./OpenSCAD/top.scad; exit"
gnome-terminal -- bash -c "openscad -o ./STLs/SSD1322_side.stl -p ./OpenSCAD/side.json -P SSD1322 ./OpenSCAD/side.scad; exit"
gnome-terminal -- bash -c "openscad -o ./STLs/SSD1322_stand_panel.stl -p ./OpenSCAD/stand_panel.json -P SSD1322 ./OpenSCAD/stand_panel.scad; exit"

gnome-terminal -- bash -c "openscad -o ./STLs/SSD1306_bottom.stl -p ./OpenSCAD/bottom.json -P SSD1306 ./OpenSCAD/bottom.scad; exit"
gnome-terminal -- bash -c "openscad -o ./STLs/SSD1306_top.stl -p ./OpenSCAD/top.json -P SSD1306 ./OpenSCAD/top.scad; exit"
gnome-terminal -- bash -c "openscad -o ./STLs/SSD1306_side.stl -p ./OpenSCAD/side.json -P SSD1306 ./OpenSCAD/side.scad; exit"
gnome-terminal -- bash -c "openscad -o ./STLs/SSD1306_stand_panel.stl -p ./OpenSCAD/stand_panel.json -P SSD1306 ./OpenSCAD/stand_panel.scad; exit"