include <BOSL2/std.scad>

// Start config
// Panel Size
panel_size = [ 104, 37 ]; // mm
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

module panel()
{
    difference()
    {
        cuboid([ panel_size.x, panel_size.y, walls_depth ], rounding = rounding_radius,
               edges = [ FRONT + RIGHT, FRONT + LEFT, BACK + RIGHT, BACK + LEFT ], $fn = resolution,
               anchor = LEFT + BOTTOM + FRONT);
        translate([ walls_thickness, walls_thickness, 0 ])
            cuboid([ panel_size.x - walls_thickness * 2, panel_size.y - walls_thickness * 2, walls_depth ],
                   rounding = rounding_radius, edges = [ FRONT + RIGHT, FRONT + LEFT, BACK + RIGHT, BACK + LEFT ],
                   $fn = resolution, anchor = LEFT + BOTTOM + FRONT);
    }
}

module raspberry_pi_pico_micro_usb_port()
{
    translate([ 0, pi_board_size.y / 2 - pi_board_usb_size.x / 2, 0 ]) cuboid(
        [ walls_thickness * 2, pi_board_usb_size.x, pi_board_usb_size.y ], rounding = pi_board_usb_rounding_radius,
        edges = [ FRONT + TOP, BACK + TOP, BOTTOM + FRONT, BOTTOM + BACK ], $fn = resolution,
        anchor = LEFT + BOTTOM + FRONT);
}

module raspberry_pi_pico()
{
    rpip_anchor = [ 0, pi_board_offset.y + screen_anchor_coords.y, 0 ];
    // translate(rpip_anchor) translate ([0,0,panel_thickness-pi_board_inset_depth])
    // cube([pi_board_size.x,pi_board_size.y,pi_board_inset_depth]);
    translate(rpip_anchor) raspberry_pi_pico_micro_usb_port();
}

screen_anchor_coords =
    [ panel_size.x / 2 - screen_pcb_size.x / 2 * screen_count, panel_size.y / 2 - active_screen_area.y / 2, 0 ];
difference()
{
    panel();
    raspberry_pi_pico();
}