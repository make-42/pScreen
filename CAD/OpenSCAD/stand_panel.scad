include <BOSL2/std.scad>
include <asanoha.scad>

// Start config
// Panel Size
panel_size = [ 106, 42 ]; // mm
// Pi Board
pi_board_size = [ 51, 21 ];                                        // mm
pi_board_screw_holes_center_from_board_edge_distance = [ 2, 4.8 ]; // mm
pi_board_screw_holes_diameter = 2.1;                               // mm
pi_board_offset = [ 1.75, -1.0 ];
pi_board_usb_size = [ 8.5, 4.5 ];
pi_board_usb_rounding_radius = 0.5; // mm
// pi_board_inset_depth = 0.2; //mm (keep at zero for final)

// Screen
screen_pcb_size = [ 100.5, 33.5, 13.3 ];                                // mm (+- 0.2) 13.3mm are pin headers.
screen_pcb_screw_holes_center_from_board_edge_distance = [ 2.25, 2.5 ]; // mm
screen_pcb_screw_holes_diameter = 2.8;                                  // mm
screen_pcb_screw_holes_spacing = [ 95, 28.5 ];                          // mm
active_screen_area = [ 76.78, 19.18 ];                                  // mm (old: 13.86)
active_screen_area_offset_from_corner = [ 11.36, 5.1 ];                 // mm
screen_rounding_radius = 0.5;                                           // mm
screen_count = 1;

// Screws
m2_screw_head_depth = 1;      // mm it's 2mm normally
m2_screw_head_diameter = 4.0; // mm (adding .6mm for tolerances)
m2_screw_nut_depth = 1;       // mm it's 1.6mm normally
m2_screw_nut_diameter = 5.0;  // mm (adding .6mm for tolerances)
m2_screw_diameter = 2.4;      // mm

// Asanoha
back_asanoha_pattern_repeat = [ 9, 3 ]; // mm
asanoha_area_rounding_radius = 3;       // mm

// Stand
stand_size = [ 3, 6, 6 ];
stand_rounding_radius = 1.5;                // mm
stand_hinge_ball_diameter = 2;              // mm
stand_hinge_ball_offset = 0;                // mm
stand_panel_tolerance = 0.25;               // mm
stand_panel_thickness = 1.5;                // mm
stand_panel_rounding_radius = 1.5;          // mm
stand_asanoha_pattern_unit_length = 10.175; // mm
stand_asanoha_repeat = [ 9, 3 ];

// Misc
walls_depth = 9.5;     // mm
walls_thickness = 1.5; // mm
panel_thickness = 1.5; // mm
rounding_radius = 1.5; // mm

resolution = 80;
// End config

stand_asanoha_pattern_size = [
    stand_asanoha_pattern_unit_length * stand_asanoha_repeat.x,
    stand_asanoha_pattern_unit_length *stand_asanoha_repeat.y
];
module asanoha_area(asanoha_anchor = [])
{
    translate(asanoha_anchor) cuboid(
        [ stand_asanoha_pattern_size.x, stand_asanoha_pattern_size.y, stand_panel_thickness ],
        rounding = asanoha_area_rounding_radius, edges = [ FRONT + RIGHT, FRONT + LEFT, BACK + RIGHT, BACK + LEFT ],
        $fn = resolution, anchor = LEFT + BOTTOM + FRONT);
}

module panel()
{
    asanoha_anchor = [
        stand_panel_size.x / 2 - stand_asanoha_pattern_size.x / 2,
        stand_panel_size.y / 2 - stand_asanoha_pattern_size.y / 2, 0
    ];
    difference()
    {
        cuboid([ stand_panel_size.x, stand_panel_size.y, stand_panel_thickness ], rounding = rounding_radius,
               edges = [ FRONT + RIGHT, FRONT + LEFT, BACK + RIGHT, BACK + LEFT ], $fn = resolution,
               anchor = LEFT + BOTTOM + FRONT);
        asanoha_area(asanoha_anchor);
    }
    intersection()
    {
        translate(asanoha_anchor) drawAsanoha(repeatXarg = stand_asanoha_repeat.x, repeatYarg = stand_asanoha_repeat.y);
        asanoha_area(asanoha_anchor);
    }
}

module hinge_ball(invert = false)
{
    stand_hinge_ball_x = invert ? stand_panel_size.x - stand_hinge_ball_offset : stand_hinge_ball_offset;
    translate([ stand_hinge_ball_x, stand_hinge_ball_diameter / 2, stand_panel_thickness / 2 ]) intersection()
    {
        sphere(d = stand_hinge_ball_diameter, $fn = resolution);
        cube([ stand_hinge_ball_diameter, stand_hinge_ball_diameter, stand_panel_thickness ], center = true);
    }
}

module hinge()
{
    hinge_ball(invert = false);
    hinge_ball(invert = true);
}

stand_panel_size = [
    panel_size.x - stand_size.x * 2 - stand_panel_tolerance * 2,
    panel_size.y - stand_size.y / 2 + stand_hinge_ball_diameter / 2
];

union()
{
    difference()
    {
        panel();
    }
    hinge();
}